package toposort_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/bees/internal/util/toposort"
)

func TestOrder_ParentBeforeChild(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: child "c" listed before its parent "p"
	ids := []string{"c", "p"}
	children := map[string][]string{"p": {"c"}}

	// Act
	got := toposort.Order(ids, children)

	// Assert: parent emitted first despite input order
	require.Equal(t, []string{"p", "c"}, got)
}

func TestOrder_MultiLevelChain(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: g -> p -> c, fed leaf-first
	ids := []string{"c", "p", "g"}
	children := map[string][]string{"g": {"p"}, "p": {"c"}}

	// Act
	got := toposort.Order(ids, children)

	// Assert: full chain ordered top-down
	require.Equal(t, []string{"g", "p", "c"}, got)
}

func TestOrder_CycleEmitsAllWithoutDropping(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: a -> b -> a (cycle); neither is a root
	ids := []string{"a", "b"}
	children := map[string][]string{"a": {"b"}, "b": {"a"}}

	// Act
	got := toposort.Order(ids, children)

	// Assert: both present exactly once — the safety net, no infinite recursion
	require.ElementsMatch(t, []string{"a", "b"}, got)
}

func TestOrder_IgnoresEdgesOutsideSet(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange: "p" points at a child not in the input set
	ids := []string{"p"}
	children := map[string][]string{"p": {"ghost"}}

	// Act
	got := toposort.Order(ids, children)

	// Assert: only the input id is returned — the unknown child isn't pulled in
	require.Equal(t, []string{"p"}, got)
}
