package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T, tmpDir, relPth string) {
	pth := filepath.Join(tmpDir, relPth)
	dirPth := filepath.Dir(pth)
	require.NoError(t, os.MkdirAll(dirPth, 0777))
	require.NoError(t, fileutil.WriteStringToFile(pth, "test"))
}

func Test_findModTimesByPath(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	for _, pth := range []string{
		"file",
		"subdir/file",
	} {
		createTestFile(t, tmpDir, pth)
	}

	modTimesByPath, err := findModTimesByPath(tmpDir)
	require.NoError(t, err)

	desiredPthMap := map[string]bool{
		filepath.Join(tmpDir, "file"):        false,
		filepath.Join(tmpDir, "subdir/file"): false,
	}
	require.Equal(t, len(desiredPthMap), len(modTimesByPath))
	for pth := range modTimesByPath {
		_, ok := desiredPthMap[pth]
		require.True(t, ok, fmt.Sprintf("%s - %#v", pth, desiredPthMap))
	}
}

func Test_isInTimeInterval(t *testing.T) {
	startTime := time.Now()
	testTime := time.Now()
	endTime := time.Now()

	require.True(t, isInTimeInterval(testTime, startTime, endTime))
	require.False(t, isInTimeInterval(testTime, endTime, startTime))
	require.False(t, isInTimeInterval(startTime, testTime, endTime))
	require.False(t, isInTimeInterval(endTime, startTime, testTime))
}

func Test_filterModTimesByPathByTimeWindow(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"file",
		"subdir/file",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "not_in_time_window")

	modTimesByPath, err := findModTimesByPath(tmpDir)
	require.NoError(t, err)

	desiredPthMap := map[string]bool{
		filepath.Join(tmpDir, "file"):               false,
		filepath.Join(tmpDir, "subdir/file"):        false,
		filepath.Join(tmpDir, "not_in_time_window"): false,
	}
	require.Equal(t, len(desiredPthMap), len(modTimesByPath))
	for pth := range modTimesByPath {
		_, ok := desiredPthMap[pth]
		require.True(t, ok, fmt.Sprintf("%s - %#v", pth, desiredPthMap))
	}

	t.Log("returns file infos of files modified within start and end time")
	{
		filteredModTimesByPath := filterModTimesByPathByTimeWindow(modTimesByPath, startTime, endTime)
		desiredPthMap := map[string]bool{
			filepath.Join(tmpDir, "file"):        false,
			filepath.Join(tmpDir, "subdir/file"): false,
		}
		require.Equal(t, len(desiredPthMap), len(filteredModTimesByPath))
		for pth := range filteredModTimesByPath {
			_, ok := desiredPthMap[pth]
			require.True(t, ok, fmt.Sprintf("%s - %#v", pth, desiredPthMap))
		}
	}

	t.Log("returns empty list if zero time window")
	{
		filteredInfos := filterModTimesByPathByTimeWindow(modTimesByPath, time.Time{}, time.Time{})
		require.Equal(t, 0, len(filteredInfos))
	}

	t.Log("returns empty list if invalid time window")
	{
		filteredInfos := filterModTimesByPathByTimeWindow(modTimesByPath, endTime, startTime)
		require.Equal(t, 0, len(filteredInfos))
	}
}

func Test_findLastModifiedPathWithFileNameRegexps(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	for _, pth := range []string{
		"subdir/apk",
		"file1",
		"subdir/file2",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	modTimesByPath, err := findModTimesByPath(tmpDir)
	require.NoError(t, err)

	desiredPthMap := map[string]bool{
		filepath.Join(tmpDir, "subdir/apk"):   false,
		filepath.Join(tmpDir, "file1"):        false,
		filepath.Join(tmpDir, "subdir/file2"): false,
	}
	require.Equal(t, len(desiredPthMap), len(modTimesByPath))
	for pth := range modTimesByPath {
		_, ok := desiredPthMap[pth]
		require.True(t, ok, fmt.Sprintf("%s - %#v", pth, desiredPthMap))
	}

	t.Log("filename regexp")
	{
		pth := findLastModifiedPathWithFileNameRegexps(modTimesByPath, regexp.MustCompile("apk"))
		require.Equal(t, filepath.Join(tmpDir, "subdir/apk"), pth)
	}

	t.Log("it finds the last modified")
	{
		pth := findLastModifiedPathWithFileNameRegexps(modTimesByPath, regexp.MustCompile("file.*"))
		require.Equal(t, filepath.Join(tmpDir, "subdir/file2"), pth)
	}

	t.Log("returns the last modified file without regex")
	{
		pth := findLastModifiedPathWithFileNameRegexps(modTimesByPath)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file2"), pth)
	}

	t.Log("returns the most strict match")
	{
		pth := findLastModifiedPathWithFileNameRegexps(modTimesByPath, regexp.MustCompile("file1"), regexp.MustCompile("file.*"))
		require.Equal(t, filepath.Join(tmpDir, "file1"), pth)
	}
}

func Test_findArtifact(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"file1",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "subdir/file2")

	pth, err := findArtifact(tmpDir, startTime, endTime, "file.*")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(tmpDir, "file1"), pth)
}

func Test_exportApk(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"file-signed.apk",
		"file.apk",
		"artifact.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "artifact-signed.apk")

	t.Log("time window test")
	{
		pth, err := exportApk(tmpDir, "artifact", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.apk"), pth)
	}

	t.Log("it prefres signed apk")
	{
		pth, err := exportApk(tmpDir, "file", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "file-signed.apk"), pth)
	}

	t.Log("it returns latest signed apk if artificat name does not match")
	{
		pth, err := exportApk(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "file-signed.apk"), pth)
	}
}

func Test_exportIpa(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"subdir/file.ipa",
		"artifact.ipa",
		"file.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "file1.ipa")

	t.Log("time window test")
	{
		pth, err := exportIpa(tmpDir, "file*", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file.ipa"), pth)
	}

	t.Log("it returns latest ipa if artificat name does not match")
	{
		pth, err := exportIpa(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.ipa"), pth)
	}
}

func Test_exportXCArchive(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"subdir/file.xcarchive",
		"artifact.xcarchive",
		"file.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "file1.xcarchive")

	t.Log("time window test")
	{
		pth, err := exportXCArchive(tmpDir, "file*", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file.xcarchive"), pth)
	}

	t.Log("it returns latest xcarchive if artificat name does not match")
	{
		pth, err := exportXCArchive(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.xcarchive"), pth)
	}
}

func Test_exportAppDSYM(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"subdir/file.app.dSYM",
		"artifact.app.dSYM",
		"file.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "file1.app.dSYM")

	t.Log("time window test")
	{
		pth, err := exportAppDSYM(tmpDir, "file*", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file.app.dSYM"), pth)
	}

	t.Log("it returns latest app.dSYM if artificat name does not match")
	{
		pth, err := exportAppDSYM(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.app.dSYM"), pth)
	}
}

func Test_exportFrameworkDSYMs(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	for _, pth := range []string{
		"artifact.framework.dSYM",
		"file.framework.dSYM",
	} {
		createTestFile(t, tmpDir, pth)
	}

	createTestFile(t, tmpDir, "file1.app.dSYM")

	pths, err := exportFrameworkDSYMs(tmpDir)
	require.NoError(t, err)
	require.Equal(t, 2, len(pths))

	desiredPthMap := map[string]bool{
		filepath.Join(tmpDir, "artifact.framework.dSYM"): false,
		filepath.Join(tmpDir, "file.framework.dSYM"):     false,
	}
	require.Equal(t, len(desiredPthMap), len(pths))
	for _, pth := range pths {
		_, ok := desiredPthMap[pth]
		require.True(t, ok, fmt.Sprintf("%s - %#v", pth, desiredPthMap))
	}
}

func Test_exportPKG(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"subdir/file.pkg",
		"artifact.pkg",
		"file.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "file1.pkg")

	t.Log("time window test")
	{
		pth, err := exportPKG(tmpDir, "file*", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file.pkg"), pth)
	}

	t.Log("it returns latest pkg if artificat name does not match")
	{
		pth, err := exportPKG(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.pkg"), pth)
	}
}

func Test_exportApp(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"subdir/file.app",
		"artifact.app",
		"file.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "file1.app")

	t.Log("time window test")
	{
		pth, err := exportApp(tmpDir, "file*", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file.app"), pth)
	}

	t.Log("it returns latest app if artificat name does not match")
	{
		pth, err := exportApp(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.app"), pth)
	}
}

func Test_exportDLL(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("file_infos_test")
	require.NoError(t, err)

	startTime := time.Now()
	time.Sleep(3 * time.Second)

	for _, pth := range []string{
		"subdir/file.dll",
		"artifact.dll",
		"file.apk",
	} {
		createTestFile(t, tmpDir, pth)
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Second)
	endTime := time.Now()
	time.Sleep(3 * time.Second)

	createTestFile(t, tmpDir, "file1.dll")

	t.Log("time window test")
	{
		pth, err := exportDLL(tmpDir, "file*", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "subdir/file.dll"), pth)
	}

	t.Log("it returns latest app if artificat name does not match")
	{
		pth, err := exportDLL(tmpDir, "does not match", startTime, endTime)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "artifact.dll"), pth)
	}
}
