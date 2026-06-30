package service

import "github.com/w-h-a/bees/v2/internal/domain"

type ContextSummary struct {
	InProgress   []domain.Issue
	Ready        []domain.Issue
	Blocked      []domain.Issue
	RecentlyDone []domain.Issue
}
