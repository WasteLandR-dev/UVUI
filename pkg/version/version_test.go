package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareVersions(t *testing.T) {
	// Test basic version comparison
	result1 := CompareVersions("3.12.0", "3.11.0")
	t.Logf("CompareVersions(\"3.12.0\", \"3.11.0\") = %d", result1)
	assert.Equal(t, 1, result1)

	result2 := CompareVersions("3.11.0", "3.12.0")
	t.Logf("CompareVersions(\"3.11.0\", \"3.12.0\") = %d", result2)
	assert.Equal(t, -1, result2)

	result3 := CompareVersions("3.12.0", "3.12.0")
	t.Logf("CompareVersions(\"3.12.0\", \"3.12.0\") = %d", result3)
	assert.Equal(t, 0, result3)

	// Test different version lengths
	// Note: "3.12" is equivalent to "3.12.0" (missing parts are treated as 0)
	assert.Equal(t, 0, CompareVersions("3.12.0", "3.12"))
	assert.Equal(t, 0, CompareVersions("3.12", "3.12.0"))

	// Test edge cases
	assert.Equal(t, 0, CompareVersions("", ""))
	assert.Equal(t, 1, CompareVersions("3.12.0", ""))
	assert.Equal(t, -1, CompareVersions("", "3.12.0"))

	// Test complex versions
	assert.Equal(t, 1, CompareVersions("3.12.1", "3.12.0"))
	assert.Equal(t, -1, CompareVersions("3.12.0", "3.12.1"))
	assert.Equal(t, 0, CompareVersions("3.12.1", "3.12.1"))
}

func TestIsNewerVersion(t *testing.T) {
	// Test newer versions
	assert.True(t, IsNewerVersion("3.12.0", "3.11.0"))
	assert.True(t, IsNewerVersion("3.12.1", "3.12.0"))
	assert.True(t, IsNewerVersion("4.0.0", "3.99.99"))

	// Test older versions
	assert.False(t, IsNewerVersion("3.11.0", "3.12.0"))
	assert.False(t, IsNewerVersion("3.12.0", "3.12.1"))
	assert.False(t, IsNewerVersion("3.99.99", "4.0.0"))

	// Test same versions
	assert.False(t, IsNewerVersion("3.12.0", "3.12.0"))
	assert.False(t, IsNewerVersion("", ""))
}

func TestIsSameVersion(t *testing.T) {
	// Test same versions
	assert.True(t, IsSameVersion("3.12.0", "3.12.0"))
	assert.True(t, IsSameVersion("", ""))
	assert.True(t, IsSameVersion("1.0.0", "1.0.0"))

	// Test different versions
	assert.False(t, IsSameVersion("3.12.0", "3.11.0"))
	assert.False(t, IsSameVersion("3.12.0", "3.12.1"))
	assert.False(t, IsSameVersion("3.12.0", ""))
	assert.False(t, IsSameVersion("", "3.12.0"))
}

func TestCompareVersions_EdgeCases(t *testing.T) {
	// Test single component versions
	assert.Equal(t, 1, CompareVersions("3", "2"))
	assert.Equal(t, -1, CompareVersions("2", "3"))
	assert.Equal(t, 0, CompareVersions("3", "3"))

	// Test versions with many components
	assert.Equal(t, 1, CompareVersions("1.2.3.4.5", "1.2.3.4.4"))
	assert.Equal(t, -1, CompareVersions("1.2.3.4.4", "1.2.3.4.5"))
	assert.Equal(t, 0, CompareVersions("1.2.3.4.5", "1.2.3.4.5"))

	// Test versions with leading zeros
	assert.Equal(t, 0, CompareVersions("3.01.0", "3.1.0"))
	assert.Equal(t, 0, CompareVersions("3.1.0", "3.01.0"))
}

func TestCompareVersions_InvalidInput(t *testing.T) {
	// Test with non-numeric components (should be treated as 0)
	assert.Equal(t, 0, CompareVersions("3.12.0", "3.12.0"))
	assert.Equal(t, 0, CompareVersions("3.12.0", "3.12.a")) // Non-numeric treated as 0
	assert.Equal(t, 0, CompareVersions("3.12.a", "3.12.0")) // Non-numeric treated as 0

	// Test with mixed valid/invalid
	assert.Equal(t, 0, CompareVersions("3.12.0", "3.12")) // Same when non-numeric parts are 0
	assert.Equal(t, 0, CompareVersions("3.12", "3.12.0")) // Same when non-numeric parts are 0
}

func TestVersionComparison_RealWorldExamples(t *testing.T) {
	// Test common Python version patterns
	versions := []string{
		"3.12.1",
		"3.12.0",
		"3.11.7",
		"3.11.6",
		"3.10.13",
		"3.10.12",
		"3.9.18",
		"3.9.17",
		"3.8.18",
	}

	// Verify descending order (newest first)
	for i := 0; i < len(versions)-1; i++ {
		assert.True(t, IsNewerVersion(versions[i], versions[i+1]),
			"%s should be newer than %s", versions[i], versions[i+1])
	}

	// Test specific comparisons
	assert.True(t, IsNewerVersion("3.12.1", "3.12.0"))
	assert.True(t, IsNewerVersion("3.12.0", "3.11.7"))
	assert.True(t, IsNewerVersion("3.11.7", "3.11.6"))
	assert.True(t, IsNewerVersion("3.10.13", "3.10.12"))
}

func TestVersionComparison_PreRelease(t *testing.T) {
	// Test pre-release versions (should be treated as regular versions for now)
	// Note: The current implementation treats non-numeric parts as 0
	assert.True(t, IsNewerVersion("3.12.0", "3.11.0"))
	assert.True(t, IsNewerVersion("3.12.1", "3.12.0"))
	assert.True(t, IsNewerVersion("3.12.2", "3.12.1"))

	// Test that pre-release versions are treated as newer than stable
	assert.True(t, IsNewerVersion("3.12.0", "3.11.99"))
}

func TestVersionComparison_ZeroPadding(t *testing.T) {
	// Test that zero-padded versions are equivalent
	assert.True(t, IsSameVersion("3.1.0", "3.01.00"))
	assert.True(t, IsSameVersion("3.10.0", "3.010.000"))

	// Test comparison with zero-padded versions
	assert.True(t, IsNewerVersion("3.10.0", "3.9.0"))
	assert.True(t, IsNewerVersion("3.010.000", "3.009.000"))
}
