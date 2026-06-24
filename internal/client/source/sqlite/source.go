package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/w-h-a/bees/internal/client/source"
	"github.com/w-h-a/bees/internal/domain"
	_ "modernc.org/sqlite"
)

// sqliteReader reads v1 bees stores. It holds no connection: each
// Read opens a fresh immutable handle and closes it becaue migrate
// reads each source exactly once.
type sqliteReader struct {
	options source.Options
}

func NewReader(opts ...source.Option) (source.Reader, error) {
	options := source.NewOptions(opts...)

	s := &sqliteReader{
		options: options,
	}

	return s, nil
}

// Read opens dbPath read-only, reads every issue with its relations, and closes
// the handle. open returns source.ErrNotFound for an absent file; Read passes it
// through unwrapped so the caller can match it and skip+report.
func (r *sqliteReader) Read(ctx context.Context, dbPath string) ([]domain.Issue, error) {
	db, err := r.open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	issues, err := r.issues(ctx, db)
	if err != nil {
		return nil, err
	}

	labels, err := r.labelsByIssue(ctx, db)
	if err != nil {
		return nil, err
	}

	deps, err := r.dependenciesByIssue(ctx, db)
	if err != nil {
		return nil, err
	}

	comments, err := r.commentsByIssue(ctx, db)
	if err != nil {
		return nil, err
	}

	handoffs, err := r.handoffsByIssue(ctx, db)
	if err != nil {
		return nil, err
	}

	for i := range issues {
		id := issues[i].ID
		issues[i].Labels = labels[id]
		issues[i].Dependencies = deps[id]
		issues[i].Comments = comments[id]
		issues[i].Handoffs = handoffs[id]
	}

	return issues, nil
}

// open prepares a handle for one source.
func (r *sqliteReader) open(dbPath string) (*sql.DB, error) {
	_, err := os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, source.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to state source %q: %w", dbPath, err)
	}

	abs, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve source path %q: %w", dbPath, err)
	}

	dsn := (&url.URL{Scheme: "file", Path: abs, RawQuery: "immutable=1"}).String()

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open source %s: %w", dbPath, err)
	}

	return db, nil
}

// issues reads every issue row from the open source handle, mirroring the repo
// adapter's column set.
func (r *sqliteReader) issues(ctx context.Context, db *sql.DB) ([]domain.Issue, error) {
	rows, err := db.QueryContext(
		ctx,
		"SELECT id, title, description, status, type, priority, assignee, estimate_mins, defer_until, due_at, created_at, updated_at, closed_at, parent_id FROM issues",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read issues: %w", err)
	}
	defer rows.Close()

	var issues []domain.Issue
	for rows.Next() {
		var issue domain.Issue
		if err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Description,
			&issue.Status,
			&issue.Type,
			&issue.Priority,
			&issue.Assignee,
			&issue.EstimateMins,
			&issue.DeferUntil,
			&issue.DueAt,
			&issue.CreatedAt,
			&issue.UpdatedAt,
			&issue.ClosedAt,
			&issue.ParentID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan issue: %w", err)
		}
		issues = append(issues, issue)
	}

	return issues, rows.Err()
}

// labelsByIssue reads all labels from the open handle, grouped by issue ID.
func (r *sqliteReader) labelsByIssue(ctx context.Context, db *sql.DB) (map[string][]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT issue_id, label FROM labels ORDER BY issue_id, label")
	if err != nil {
		return nil, fmt.Errorf("failed to read labels: %w", err)
	}
	defer rows.Close()

	byIssue := map[string][]string{}
	for rows.Next() {
		var issueID, label string
		if err := rows.Scan(&issueID, &label); err != nil {
			return nil, fmt.Errorf("failed to scan label: %w", err)
		}
		byIssue[issueID] = append(byIssue[issueID], label)
	}

	return byIssue, rows.Err()
}

// dependenciesByIssue reads all dependencies from the open handle, grouped by issue ID.
func (r *sqliteReader) dependenciesByIssue(ctx context.Context, db *sql.DB) (map[string][]domain.Dependency, error) {
	rows, err := db.QueryContext(ctx, "SELECT issue_id, depends_on_id, created_at FROM dependencies ORDER BY issue_id, created_at")
	if err != nil {
		return nil, fmt.Errorf("failed to read dependencies: %w", err)
	}
	defer rows.Close()

	byIssue := map[string][]domain.Dependency{}
	for rows.Next() {
		var d domain.Dependency
		if err := rows.Scan(&d.IssueID, &d.DependsOnID, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %w", err)
		}
		byIssue[d.IssueID] = append(byIssue[d.IssueID], d)
	}

	return byIssue, rows.Err()
}

// commentsByIssue reads all comments from the open handle, grouped by issue ID.
func (r *sqliteReader) commentsByIssue(ctx context.Context, db *sql.DB) (map[string][]domain.Comment, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, issue_id, author, body, created_at FROM comments ORDER BY issue_id, created_at")
	if err != nil {
		return nil, fmt.Errorf("failed to read comments: %w", err)
	}
	defer rows.Close()

	byIssue := map[string][]domain.Comment{}
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.IssueID, &c.Author, &c.Body, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		byIssue[c.IssueID] = append(byIssue[c.IssueID], c)
	}

	return byIssue, rows.Err()
}

// handoffsByIssue reads all handoffs from the open handle, grouped by issue ID.
func (r *sqliteReader) handoffsByIssue(ctx context.Context, db *sql.DB) (map[string][]domain.Handoff, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, issue_id, done, remaining, decisions, uncertain, created_at FROM handoffs ORDER BY issue_id, created_at")
	if err != nil {
		return nil, fmt.Errorf("failed to read handoffs: %w", err)
	}
	defer rows.Close()

	byIssue := map[string][]domain.Handoff{}
	for rows.Next() {
		var h domain.Handoff
		if err := rows.Scan(&h.ID, &h.IssueID, &h.Done, &h.Remaining, &h.Decisions, &h.Uncertain, &h.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan handoff: %w", err)
		}
		byIssue[h.IssueID] = append(byIssue[h.IssueID], h)
	}

	return byIssue, rows.Err()
}
