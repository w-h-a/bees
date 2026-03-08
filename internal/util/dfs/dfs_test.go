package dfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectCycle_EmptyGraph(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Act
	hasCycle, cycle := DetectCycle(nil, "a")

	// Assert
	require.False(t, hasCycle)
	require.Empty(t, cycle)
}

func TestDetectCycle_LinearChain(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b"},
		"b": {"c"},
	}

	// Act
	hasCycle, cycle := DetectCycle(graph, "a")

	// Assert
	require.False(t, hasCycle)
	require.Empty(t, cycle)
}

func TestDetectCycle_Diamond(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
	}

	// Act
	hasCycle, cycle := DetectCycle(graph, "a")

	// Assert
	require.False(t, hasCycle)
	require.Empty(t, cycle)
}

func TestDetectCycle_TwoNode(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}

	// Act
	hasCycle, cycle := DetectCycle(graph, "a")

	// Assert
	require.True(t, hasCycle)
	require.Equal(t, []string{"a", "b", "a"}, cycle)
}

func TestDetectCycle_ThreeNode(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"a"},
	}

	// Act
	hasCycle, cycle := DetectCycle(graph, "a")

	// Assert
	require.True(t, hasCycle)
	require.Equal(t, []string{"a", "b", "c", "a"}, cycle)
}

func TestDetectCycle_CycleNotInvolvingStart(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b"},
		"b": {"c"},
		"c": {"b"},
	}

	// Act
	hasCycle, cycle := DetectCycle(graph, "a")

	// Assert
	require.True(t, hasCycle)
	require.Equal(t, []string{"b", "c", "b"}, cycle)
}

func TestDetectCycle_StartNotInGraph(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"x": {"y"},
	}

	// Act
	hasCycle, cycle := DetectCycle(graph, "a")

	// Assert
	require.False(t, hasCycle)
	require.Empty(t, cycle)
}

func TestReachable_EmptyGraph(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Act
	result := Reachable(nil, "a")

	// Assert
	require.Equal(t, map[string]bool{"a": true}, result)
}

func TestReachable_LinearChain(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b"},
		"b": {"c"},
	}

	// Act
	result := Reachable(graph, "a")

	// Assert
	require.Equal(t, map[string]bool{"a": true, "b": true, "c": true}, result)
}

func TestReachable_StartNotInGraph(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"x": {"y"},
	}

	// Act
	result := Reachable(graph, "a")

	// Assert
	require.Equal(t, map[string]bool{"a": true}, result)
}

func TestReachable_Diamond(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	graph := map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
	}

	// Act
	result := Reachable(graph, "a")

	// Assert
	require.Equal(t, map[string]bool{"a": true, "b": true, "c": true, "d": true}, result)
}
