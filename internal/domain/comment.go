package domain

import "time"

type Comment struct {
	ID        int64
	IssueID   string
	Author    string
	Body      string
	CreatedAt time.Time
}
