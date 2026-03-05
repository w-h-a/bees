package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/w-h-a/bees/internal/domain"
)

func issueFields(issue *domain.Issue) map[string]string {
	fields := map[string]string{
		"title":         issue.Title,
		"description":   issue.Description,
		"status":        string(issue.Status),
		"type":          string(issue.Type),
		"assignee":      issue.Assignee,
		"estimate_mins": fmt.Sprintf("%d", issue.EstimateMins),
	}
	if issue.Priority != nil {
		fields["priority"] = fmt.Sprintf("%d", *issue.Priority)
	}
	if issue.DeferUntil != nil {
		fields["defer_until"] = issue.DeferUntil.UTC().Format(time.RFC3339)
	}
	if issue.DueAt != nil {
		fields["due_at"] = issue.DueAt.UTC().Format(time.RFC3339)
	}
	if issue.ClosedAt != nil {
		fields["closed_at"] = issue.ClosedAt.UTC().Format(time.RFC3339)
	}
	if issue.ParentID != nil {
		fields["parent_id"] = *issue.ParentID
	}
	labels := make([]string, len(issue.Labels))
	copy(labels, issue.Labels)
	sort.Strings(labels)
	for i, l := range labels {
		fields[fmt.Sprintf("label_%d", i)] = l
	}
	return fields
}
