package beads

import (
	"time"

	"github.com/w-h-a/bees/internal/domain"
)

func mapToIssue(bi beadsIssue) domain.Issue {
	issue := domain.Issue{
		ID:           bi.ID,
		Title:        bi.Title,
		Description:  bi.Description,
		Status:       mapStatus(bi.Status),
		Type:         domain.Type(bi.IssueType),
		Priority:     bi.Priority,
		Assignee:     bi.Assignee,
		EstimateMins: bi.EstimatedMinutes,
		DeferUntil:   bi.DeferUntil,
		CreatedAt:    bi.CreatedAt,
		UpdatedAt:    bi.UpdatedAt,
		ClosedAt:     bi.ClosedAt,
		Labels:       bi.Labels,
	}

	if bi.Status == "deferred" && bi.DeferUntil == nil {
		now := time.Now()
		issue.DeferUntil = &now
	}

	for _, d := range bi.Dependencies {
		switch d.Type {
		case "blocks":
			issue.Dependencies = append(issue.Dependencies, domain.Dependency{
				IssueID:     d.IssueID,
				DependsOnID: d.DependsOnID,
				CreatedAt:   d.CreatedAt,
			})
		case "parent-child":
			parentID := d.DependsOnID
			issue.ParentID = &parentID
		}
	}

	for _, c := range bi.Comments {
		issue.Comments = append(issue.Comments, domain.Comment{
			IssueID:   c.IssueID,
			Author:    c.Author,
			Body:      c.Text,
			CreatedAt: c.CreatedAt,
		})
	}

	if bi.CloseReason != "" {
		issue.Comments = append(issue.Comments, domain.Comment{
			IssueID:   bi.ID,
			Author:    "migration",
			Body:      "Close reason: " + bi.CloseReason,
			CreatedAt: bi.UpdatedAt,
		})
	}

	return issue
}

func mapStatus(status string) domain.Status {
	switch status {
	case "open":
		return domain.StatusOpen
	case "in_progress":
		return domain.StatusInProgress
	case "closed":
		return domain.StatusClosed
	case "deferred", "blocked", "pinned", "hooked":
		return domain.StatusOpen
	default:
		return domain.Status(status)
	}
}
