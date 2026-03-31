package domain

import "time"

type Handoff struct {
	ID        int64     `json:"id"`
	IssueID   string    `json:"issue_id"`
	Done      string    `json:"done"`
	Remaining string    `json:"remaining"`
	Decisions string    `json:"decisions"`
	Uncertain string    `json:"uncertain"`
	CreatedAt time.Time `json:"created_at"`
}
