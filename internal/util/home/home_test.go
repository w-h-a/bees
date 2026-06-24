package home_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/bees/internal/util/home"
)

func TestResolve_BeesHomeOverride(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	beesHome := "/custom/bees"
	userHome := "/home/wes"

	// Act
	resolved, err := home.Resolve(beesHome, userHome)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "/custom/bees", resolved.Home)
	require.Equal(t, "/custom/bees/bees.db", resolved.DB)
}

func TestResolve_BeesHomeOnly(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	beesHome := "/custom/bees"
	userHome := ""

	// Act
	resolved, err := home.Resolve(beesHome, userHome)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "/custom/bees", resolved.Home)
	require.Equal(t, "/custom/bees/bees.db", resolved.DB)
}

func TestResolve_FallsBackToUserHome(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	beesHome := ""
	userHome := "/home/wes"

	// Act
	resolved, err := home.Resolve(beesHome, userHome)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "/home/wes/.bees", resolved.Home)
	require.Equal(t, "/home/wes/.bees/bees.db", resolved.DB)
}

func TestResolve_BothEmptyErrors(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	beesHome := ""
	userHome := ""

	// Act
	_, err := home.Resolve(beesHome, userHome)

	// Assert
	require.EqualError(t, err, "refused to resolve bees home: BEES_HOME and user home are both empty")
}
