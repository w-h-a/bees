package domain_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/bees/internal/domain"
)

func TestPlanMigration_CountsOneSource(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	sources := []domain.SourceIssues{
		{
			Path: "/repo/a",
			Issues: []domain.Issue{
				{
					ID:           "a-0001",
					Labels:       []string{"backend", "v1"},
					Comments:     []domain.Comment{{Body: "one"}, {Body: "two"}},
					Handoffs:     []domain.Handoff{{Done: "x"}},
					Dependencies: []domain.Dependency{{IssueID: "a-0001", DependsOnID: "a-0002"}},
				},
				{ID: "a-0002"},
			},
		},
	}

	// Act
	plan := domain.PlanMigration(sources, nil)

	// Assert
	require.Len(t, plan.Sources, 1)
	require.Equal(t, domain.SourceCounts{
		Path:     "/repo/a",
		Issues:   2,
		Comments: 2,
		Handoffs: 1,
		Deps:     1,
		Labels:   2,
	}, plan.Sources[0])
	require.Empty(t, plan.Collisions)
}

func TestPlanMigration_DetectsSourceVsSourceCollision(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: two sources share the full ID "dup-0001"
	sources := []domain.SourceIssues{
		{Path: "/repo/a", Issues: []domain.Issue{{ID: "dup-0001"}, {ID: "a-0002"}}},
		{Path: "/repo/b", Issues: []domain.Issue{{ID: "dup-0001"}, {ID: "b-0003"}}},
	}

	// Act
	plan := domain.PlanMigration(sources, nil)

	// Assert
	require.Equal(t, []string{"dup-0001"}, plan.Collisions)
}

func TestPlanMigration_DetectsSourceVsTargetCollision(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: the target already holds "x-0001"; a source carries the same ID
	sources := []domain.SourceIssues{
		{Path: "/repo/a", Issues: []domain.Issue{{ID: "x-0001"}, {ID: "a-0002"}}},
	}
	targetIDs := []string{"x-0001", "t-0009"}

	// Act
	plan := domain.PlanMigration(sources, targetIDs)

	// Assert
	require.Equal(t, []string{"x-0001"}, plan.Collisions)
}

func TestPlanMigration_CollisionsSorted(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: two collisions, fed in reverse order across both sources
	sources := []domain.SourceIssues{
		{Path: "/repo/a", Issues: []domain.Issue{{ID: "m-0002"}, {ID: "m-0001"}}},
		{Path: "/repo/b", Issues: []domain.Issue{{ID: "m-0002"}, {ID: "m-0001"}}},
	}

	// Act
	plan := domain.PlanMigration(sources, nil)

	// Assert: deterministic sorted order despite random map iteration
	require.Equal(t, []string{"m-0001", "m-0002"}, plan.Collisions)
}

func TestPlanMigration_EmptySourceStillReported(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: a source that exists but holds no issues
	sources := []domain.SourceIssues{
		{Path: "/repo/empty", Issues: nil},
	}

	// Act
	plan := domain.PlanMigration(sources, nil)

	// Assert: reported with an all-zero tally, not dropped
	require.Len(t, plan.Sources, 1)
	require.Equal(t, domain.SourceCounts{Path: "/repo/empty"}, plan.Sources[0])
	require.Empty(t, plan.Collisions)
}

func TestPlanMigration_MultipleSourcesCountedIndependently(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: two distinct sources, no shared IDs
	sources := []domain.SourceIssues{
		{
			Path:   "/repo/a",
			Issues: []domain.Issue{{ID: "a-0001", Labels: []string{"x"}}},
		},
		{
			Path: "/repo/b",
			Issues: []domain.Issue{
				{ID: "b-0001", Comments: []domain.Comment{{Body: "c"}}},
				{ID: "b-0002"},
			},
		},
	}

	// Act
	plan := domain.PlanMigration(sources, nil)

	// Assert: per-source tallies, in input order, no cross-source bleed
	require.Equal(t, []domain.SourceCounts{
		{Path: "/repo/a", Issues: 1, Labels: 1},
		{Path: "/repo/b", Issues: 2, Comments: 1},
	}, plan.Sources)
	require.Empty(t, plan.Collisions)
}
