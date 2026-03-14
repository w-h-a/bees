package duration

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParse_RelativeDays(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	frozen := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozen }
	t.Cleanup(func() { timeNow = time.Now })

	// Act
	result, err := Parse("1d")

	// Assert
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 3, 13, 0, 0, 0, 0, time.UTC), result)
}

func TestParse_RelativeWeeks(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	frozen := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozen }
	t.Cleanup(func() { timeNow = time.Now })

	// Act
	result, err := Parse("2w")

	// Assert
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC), result)
}

func TestParse_RelativeMonths(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	frozen := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozen }
	t.Cleanup(func() { timeNow = time.Now })

	// Act
	result, err := Parse("6mo")

	// Assert
	require.NoError(t, err)
	require.Equal(t, time.Date(2025, 9, 14, 0, 0, 0, 0, time.UTC), result)
}

func TestParse_RelativeYears(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	frozen := time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return frozen }
	t.Cleanup(func() { timeNow = time.Now })

	// Act
	result, err := Parse("1y")

	// Assert
	require.NoError(t, err)
	require.Equal(t, time.Date(2025, 3, 14, 0, 0, 0, 0, time.UTC), result)
}

func TestParse_AbsoluteDate(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Act
	result, err := Parse("2025-09-14")

	// Assert
	require.NoError(t, err)
	require.Equal(t, time.Date(2025, 9, 14, 0, 0, 0, 0, time.UTC), result)
}

func TestParse_InvalidString(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Act
	_, err := Parse("abc")

	// Assert
	require.EqualError(t, err, "invalid duration: abc")
}

func TestParse_EmptyString(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Act
	_, err := Parse("")

	// Assert
	require.EqualError(t, err, "empty duration string")
}
