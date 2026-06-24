package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/w-h-a/bees/internal/domain"
	"github.com/w-h-a/bees/internal/service"
)

// Handler is the CLI handler for bees commands. It holds the service and the
// resolved output mode, and exposes one method per command. Methods write to
// an io.Writer so they can be exercised with a buffer in tests.
type Handler struct {
	svc  *service.Service
	json bool
}

func New(svc *service.Service, jsonOut bool) *Handler {
	return &Handler{svc: svc, json: jsonOut}
}

// Version writes version information to out.
func (h *Handler) Version(out io.Writer) error {
	v, c, d := "dev", "unknown", "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				c = s.Value
			case "vcs.time":
				d = s.Value
			}
		}
	}

	if !h.json {
		short := c
		if len(short) > 7 {
			short = short[:7]
		}
		fmt.Fprintf(out, "bees %s (commit: %s, built: %s)\n", v, short, d)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(map[string]string{
		"version": v,
		"commit":  c,
		"date":    d,
	})
}

// Close closes the issue with the given id and writes the result to out.
func (h *Handler) Close(ctx context.Context, out io.Writer, id string) error {
	issue, changed, err := h.svc.CloseIssue(ctx, id)
	if err != nil {
		return err
	}

	if !h.json {
		if !changed {
			fmt.Fprintf(out, "Already closed: %s\n", issue.ID)
		} else {
			fmt.Fprintf(out, "Closed %s: %s\n", issue.ID, issue.Title)
		}
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issue)
}

// Reopen reopens the issue with the given id and writes the result to out.
func (h *Handler) Reopen(ctx context.Context, out io.Writer, id string) error {
	issue, changed, err := h.svc.ReopenIssue(ctx, id)
	if err != nil {
		return err
	}

	if !h.json {
		if !changed {
			fmt.Fprintf(out, "Already open: %s\n", issue.ID)
		} else {
			fmt.Fprintf(out, "Reopened %s: %s\n", issue.ID, issue.Title)
		}
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issue)
}

// Comment adds a comment to an issue and writes the result to out. The author
// is resolved by the caller (it may come from config).
func (h *Handler) Comment(ctx context.Context, out io.Writer, id, author, text string) error {
	comment, err := h.svc.AddComment(ctx, id, author, text)
	if err != nil {
		return err
	}

	if !h.json {
		name := comment.Author
		if name == "" {
			name = "anonymous"
		}
		ts := comment.CreatedAt.Format("2006-01-02 15:04")
		fmt.Fprintf(out, "%s %s\n", dimStyle.Render(ts), headerStyle.Render(name))
		for line := range strings.SplitSeq(comment.Body, "\n") {
			fmt.Fprintf(out, "  %s\n", line)
		}
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(comment)
}

// Handoff records a structured handoff for an issue and writes the result to out.
func (h *Handler) Handoff(ctx context.Context, out io.Writer, id, done, remaining, decisions, uncertain string) error {
	handoff, err := h.svc.AddHandoff(ctx, id, done, remaining, decisions, uncertain)
	if err != nil {
		return err
	}

	if !h.json {
		ts := handoff.CreatedAt.Format("2006-01-02 15:04")
		fmt.Fprintf(out, "%s %s\n", dimStyle.Render(ts), headerStyle.Render("Handoff recorded"))
		printHandoffInline(out, *handoff)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(handoff)
}

// Create creates an issue (already constructed and validated by the caller) and
// writes the result to out.
func (h *Handler) Create(ctx context.Context, out io.Writer, issue *domain.Issue) error {
	id, err := h.svc.CreateIssue(ctx, issue)
	if err != nil {
		return err
	}

	if !h.json {
		fmt.Fprintf(out, "Created %s: %s\n", id, issue.Title)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	result := map[string]any{
		"id":     id,
		"title":  issue.Title,
		"type":   string(issue.Type),
		"status": string(issue.Status),
	}
	if issue.Priority != nil {
		result["priority"] = *issue.Priority
	}

	return enc.Encode(result)
}

// Update applies an update to an issue and writes the result to out.
func (h *Handler) Update(ctx context.Context, out io.Writer, id string, update domain.IssueUpdate) error {
	issue, err := h.svc.UpdateIssue(ctx, id, update)
	if err != nil {
		return err
	}

	if !h.json {
		fmt.Fprintf(out, "Updated %s: %s\n", issue.ID, issue.Title)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issue)
}

// Delete previews or performs a bulk delete and writes the result to out. When
// confirm is false it only previews the candidates.
func (h *Handler) Delete(ctx context.Context, out io.Writer, filter domain.DeleteFilter, confirm bool) error {
	if !confirm {
		candidates, err := h.svc.PreviewDeleteIssues(ctx, filter)
		if err != nil {
			return err
		}

		if h.json {
			enc := json.NewEncoder(out)
			enc.SetIndent("", " ")
			return enc.Encode(candidates)
		}

		for i, c := range candidates {
			if i >= 20 {
				fmt.Fprintf(out, "... and %d more\n", len(candidates)-20)
				break
			}
			title := c.Title
			if len(title) > 50 {
				title = title[:47] + "..."
			}
			closedAt := ""
			if c.ClosedAt != nil {
				closedAt = c.ClosedAt.Format("2006-01-02")
			}
			fmt.Fprintf(out, "  %s  %-12s  %s\n", c.ID, closedAt, title)
		}

		fmt.Fprintf(out, "Would delete %d issues closed before %s. Run with --yes to confirm.\n",
			len(candidates), filter.ClosedBefore.Format("2006-01-02"))

		return nil
	}

	count, err := h.svc.DeleteIssues(ctx, filter)
	if err != nil {
		return err
	}

	if h.json {
		enc := json.NewEncoder(out)
		enc.SetIndent("", " ")
		return enc.Encode(map[string]any{"deleted": count})
	}

	fmt.Fprintf(out, "Deleted %d issues closed before %s.\n",
		count, filter.ClosedBefore.Format("2006-01-02"))

	return nil
}

// Show writes the details of a single issue to out.
func (h *Handler) Show(ctx context.Context, out io.Writer, id string) error {
	issue, err := h.svc.GetIssue(ctx, id)
	if err != nil {
		return err
	}

	if !h.json {
		printIssue(out, issue)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issue)
}

// List writes the issues matching filter to out.
func (h *Handler) List(ctx context.Context, out io.Writer, filter domain.ListFilter) error {
	issues, err := h.svc.ListIssues(ctx, filter)
	if err != nil {
		return err
	}

	if !h.json {
		printIssueTable(out, issues)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issues)
}

// Ready writes the issues ready to work on to out.
func (h *Handler) Ready(ctx context.Context, out io.Writer, sort string, limit int) error {
	issues, err := h.svc.ReadyIssues(ctx, sort, limit)
	if err != nil {
		return err
	}

	if !h.json {
		printIssueTable(out, issues)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issues)
}

// Upcoming writes the issues scheduled for the coming days to out.
func (h *Handler) Upcoming(ctx context.Context, out io.Writer, days int, assignee string) error {
	issues, err := h.svc.UpcomingIssues(ctx, days, assignee)
	if err != nil {
		return err
	}

	if !h.json {
		printUpcomingTable(out, issues)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issues)
}

// Context writes the current context summary to out.
func (h *Handler) Context(ctx context.Context, out io.Writer) error {
	summary, err := h.svc.Context(ctx)
	if err != nil {
		return err
	}

	if !h.json {
		printContextSummary(out, summary)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(summary)
}

// Search writes issues matching query to out.
func (h *Handler) Search(ctx context.Context, out io.Writer, query string, limit int) error {
	issues, err := h.svc.SearchIssues(ctx, query, limit)
	if err != nil {
		return err
	}

	if !h.json {
		printIssueTable(out, issues)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(issues)
}

// DepAdd adds a blocking dependency and writes the result to out.
func (h *Handler) DepAdd(ctx context.Context, out io.Writer, blockerID, blockedID string) error {
	blocker, blocked, err := h.svc.AddDependency(ctx, blockerID, blockedID)
	if err != nil {
		return err
	}

	if !h.json {
		fmt.Fprintf(out, "%s now blocks %s\n", blocker, blocked)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(map[string]string{
		"blocker_id": blocker,
		"blocked_id": blocked,
		"action":     "added",
	})
}

// DepRemove removes a blocking dependency and writes the result to out.
func (h *Handler) DepRemove(ctx context.Context, out io.Writer, blockerID, blockedID string) error {
	blocker, blocked, changed, err := h.svc.RemoveDependency(ctx, blockerID, blockedID)
	if err != nil {
		return err
	}

	if !h.json {
		if !changed {
			fmt.Fprintf(out, "No dependency: %s does not block %s\n", blocker, blocked)
		} else {
			fmt.Fprintf(out, "%s no longer blocks %s\n", blocker, blocked)
		}
		return nil
	}

	action := "removed"
	if !changed {
		action = "none"
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(map[string]string{
		"blocker_id": blocker,
		"blocked_id": blocked,
		"action":     action,
	})
}

// DepGraph writes the dependency graph to out.
func (h *Handler) DepGraph(ctx context.Context, out io.Writer, id *string, status string) error {
	graph, err := h.svc.BuildGraph(ctx, id, status)
	if err != nil {
		return err
	}

	if !h.json {
		printGraph(out, graph)
		return nil
	}

	type jsonNode struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		Status       string `json:"status"`
		Priority     int    `json:"priority"`
		Type         string `json:"type"`
		DeferUntil   string `json:"defer_until,omitempty"`
		EstimateMins int    `json:"estimate_mins,omitempty"`
	}
	type jsonEdge struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	type jsonGraph struct {
		Nodes []jsonNode `json:"nodes"`
		Edges []jsonEdge `json:"edges"`
	}

	nodes := make([]jsonNode, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		deferStr := ""
		if n.DeferUntil != nil {
			deferStr = n.DeferUntil.Format("2006-01-02")
		}
		nodes = append(nodes, jsonNode{
			ID:           n.ID,
			Title:        n.Title,
			Status:       string(n.Status),
			Priority:     n.Priority,
			Type:         string(n.Type),
			DeferUntil:   deferStr,
			EstimateMins: n.EstimateMins,
		})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })

	edges := make([]jsonEdge, 0, len(graph.Edges))
	for _, e := range graph.Edges {
		edges = append(edges, jsonEdge{From: e.From, To: e.To})
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(jsonGraph{Nodes: nodes, Edges: edges})
}

// Import imports issues from r and writes the result to out.
func (h *Handler) Import(ctx context.Context, out io.Writer, r io.Reader) error {
	result, err := h.svc.ImportIssues(ctx, r)
	if err != nil {
		return err
	}

	if !h.json {
		fmt.Fprintf(out, "Imported: %d created, %d updated, %d unchanged", result.Created, result.Updated, result.Unchanged)
		if result.Skipped > 0 {
			fmt.Fprintf(out, ", %d skipped", result.Skipped)
		}
		fmt.Fprintln(out)
		return nil
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", " ")

	return enc.Encode(map[string]any{
		"created":   result.Created,
		"updated":   result.Updated,
		"unchanged": result.Unchanged,
		"skipped":   result.Skipped,
	})
}

// Export streams matching issues to dest. When outputPath is non-empty a
// confirmation is written to out.
func (h *Handler) Export(ctx context.Context, out io.Writer, dest io.Writer, filter domain.ExportFilter, outputPath string) error {
	if err := h.svc.ExportIssues(ctx, dest, filter); err != nil {
		return err
	}

	if outputPath != "" {
		fmt.Fprintf(out, "Exported to %s\n", outputPath)
	}

	return nil
}

// Migrate renders the dry-run migration plan to out.
func (h *Handler) Migrate(ctx context.Context, out io.Writer, targetDBPath string, sourceRepoPaths []string) error {
	report, err := h.svc.PlanMigration(ctx, targetDBPath, sourceRepoPaths)
	if err != nil {
		return err
	}

	if h.json {
		enc := json.NewEncoder(out)
		enc.SetIndent("", " ")
		if err := enc.Encode(report); err != nil {
			return err
		}
	} else {
		printMigrateReport(out, report)
	}

	if len(report.Collisions) > 0 {
		return fmt.Errorf("refusing to migrate: %d colliding id(s): %s",
			len(report.Collisions), strings.Join(report.Collisions, ", "))
	}

	return nil
}
