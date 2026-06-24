package domain

import "sort"

// SourceIssues is one v1 source store's fully-read contents, handed to the
// planner as plain values.
type SourceIssues struct {
	Path   string
	Issues []Issue
}

// MigratePlan is the dry-run result.
type MigratePlan struct {
	Sources    []SourceCounts
	Collisions []string
}

// SourceCounts is the per-source dry-run count. It is what migrate would import
// from one source if you ran it out of dry-run mode.
type SourceCounts struct {
	Path     string
	Issues   int
	Comments int
	Handoffs int
	Deps     int
	Labels   int
}

// PlanMigration computes the dry-run plan.
func PlanMigration(sources []SourceIssues, targetIDs []string) MigratePlan {
	plan := MigratePlan{
		Sources: make([]SourceCounts, 0, len(sources)),
	}

	idCounts := map[string]int{}

	for _, id := range targetIDs {
		idCounts[id]++
	}

	for _, src := range sources {
		sourceCounts := SourceCounts{Path: src.Path, Issues: len(src.Issues)}

		for _, issue := range src.Issues {
			sourceCounts.Comments += len(issue.Comments)
			sourceCounts.Handoffs += len(issue.Handoffs)
			sourceCounts.Deps += len(issue.Dependencies)
			sourceCounts.Labels += len(issue.Labels)
			idCounts[issue.ID]++
		}

		plan.Sources = append(plan.Sources, sourceCounts)
	}

	for id, n := range idCounts {
		if n > 1 {
			plan.Collisions = append(plan.Collisions, id)
		}
	}

	sort.Strings(plan.Collisions)

	return plan
}
