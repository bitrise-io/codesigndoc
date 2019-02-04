package utility

import (
	"os"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestNewSortablePath(t *testing.T) {
	t.Log("rel path")
	{
		expectedAbsPth, err := pathutil.AbsPath("test")
		require.NoError(t, err)

		expectedComponents := strings.Split(expectedAbsPth, string(os.PathSeparator))
		fixedExpectedComponents := []string{}
		for _, c := range expectedComponents {
			if c != "" {
				fixedExpectedComponents = append(fixedExpectedComponents, c)
			}
		}

		sortable, err := NewSortablePath("test")
		require.NoError(t, err)
		require.Equal(t, "test", sortable.Pth)
		require.Equal(t, expectedAbsPth, sortable.AbsPth)
		require.Equal(t, fixedExpectedComponents, sortable.Components)
	}

	t.Log("rel path")
	{
		expectedAbsPth, err := pathutil.AbsPath("./test")
		require.NoError(t, err)

		expectedComponents := strings.Split(expectedAbsPth, string(os.PathSeparator))
		fixedExpectedComponents := []string{}
		for _, c := range expectedComponents {
			if c != "" {
				fixedExpectedComponents = append(fixedExpectedComponents, c)
			}
		}

		sortable, err := NewSortablePath("./test")
		require.NoError(t, err)
		require.Equal(t, "./test", sortable.Pth)
		require.Equal(t, expectedAbsPth, sortable.AbsPth)
		require.Equal(t, fixedExpectedComponents, sortable.Components)
	}

	t.Log("abs path")
	{
		expectedAbsPth := "/Users/bitrise/test"
		expectedComponents := []string{"Users", "bitrise", "test"}

		sortable, err := NewSortablePath("/Users/bitrise/test")
		require.NoError(t, err)
		require.Equal(t, "/Users/bitrise/test", sortable.Pth)
		require.Equal(t, expectedAbsPth, sortable.AbsPth)
		require.Equal(t, expectedComponents, sortable.Components)
	}
}

func TestSortPathsByComponents(t *testing.T) {
	t.Log("abs paths")
	{
		paths := []string{
			"/Users/bitrise/test/test/test",
			"/Users/bitrise/test/test",
			"/Users/vagrant",
			"/Users/bitrise",
		}

		expectedSorted := []string{
			"/Users/bitrise",
			"/Users/vagrant",
			"/Users/bitrise/test/test",
			"/Users/bitrise/test/test/test",
		}
		actualSorted, err := SortPathsByComponents(paths)
		require.NoError(t, err)
		require.Equal(t, expectedSorted, actualSorted)
	}

	t.Log("rel paths")
	{
		paths := []string{
			"bitrise/test/test/test",
			"bitrise/test/test",
			"vagrant",
			"bitrise",
		}

		expectedSorted := []string{
			"bitrise",
			"vagrant",
			"bitrise/test/test",
			"bitrise/test/test/test",
		}
		actualSorted, err := SortPathsByComponents(paths)
		require.NoError(t, err)
		require.Equal(t, expectedSorted, actualSorted)
	}
}
