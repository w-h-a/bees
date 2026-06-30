package domain

import (
	"sort"

	"github.com/w-h-a/bees/v2/internal/util/toposort"
)

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

// CommitPlan is one commit run broken down per source, plus any Collisions.
type CommitPlan struct {
	Sources    []SourceImport
	Collisions []string
}

// SourceImport is one source's commit breakdown: Imports are the issues whose
// IDs are new to the global store, Skipped counts the ones whose IDs are
// already there.
type SourceImport struct {
	Path    string
	Imports []Issue
	Skipped int
}

// PlanCommit splits each source's issues into those new to the target (to write)
// and those already there (skipped by ID), and collects any full ID held by two
// or more sources.
func PlanCommit(sources []SourceIssues, targetIDs []string) CommitPlan {
	plan := CommitPlan{
		Sources: make([]SourceImport, 0, len(sources)),
	}

	inTarget := make(map[string]bool, len(targetIDs))
	for _, id := range targetIDs {
		inTarget[id] = true
	}

	idCounts := map[string]int{}
	for _, src := range sources {
		for _, issue := range src.Issues {
			idCounts[issue.ID]++
		}
	}

	for _, src := range sources {
		imp := SourceImport{Path: src.Path}

		for _, issue := range src.Issues {
			if inTarget[issue.ID] {
				imp.Skipped++
				continue
			}
			imp.Imports = append(imp.Imports, issue)
		}

		imp.Imports = orderParentsFirst(imp.Imports)

		plan.Sources = append(plan.Sources, imp)
	}

	for id, n := range idCounts {
		if n > 1 {
			plan.Collisions = append(plan.Collisions, id)
		}
	}

	sort.Strings(plan.Collisions)

	return plan
}

// orderParentsFirst returns issues with every parent ahead of its children, so a
// commit inserts them without tripping the parent_id foreign key.
func orderParentsFirst(issues []Issue) []Issue {
	byID := make(map[string]Issue, len(issues))
	ids := make([]string, 0, len(issues))
	for _, issue := range issues {
		byID[issue.ID] = issue
		ids = append(ids, issue.ID)
	}

	children := map[string][]string{}
	for _, issue := range issues {
		if issue.ParentID == nil {
			continue
		}
		children[*issue.ParentID] = append(children[*issue.ParentID], issue.ID)
	}

	ordered := make([]Issue, 0, len(issues))
	for _, id := range toposort.Order(ids, children) {
		ordered = append(ordered, byID[id])
	}

	return ordered
}
