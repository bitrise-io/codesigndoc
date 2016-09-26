package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	confExportOutputDirPath = "./codesigndoc_exports"
)

func printFinishedWithError(toolName, format string, args ...interface{}) error {
	log.Errorf(colorstring.Red("Error: ")+format, args...)
	fmt.Println()
	fmt.Println("------------------------------")
	fmt.Println("First of all " + colorstring.Red("please make sure that you can Archive your app from "+toolName+"."))
	fmt.Println("codesigndoc only works if you can archive your app from " + toolName + ".")
	fmt.Println("If you can, and you get a valid IPA file if you export from " + toolName + ",")
	fmt.Println(colorstring.Red("please create an issue") + " on GitHub at: https://github.com/bitrise-tools/codesigndoc/issues")
	fmt.Println("with as many details & logs as you can share!")
	fmt.Println("------------------------------")
	fmt.Println()

	return fmt.Errorf(format, args)
}

func printFinished() {
	fmt.Println()
	fmt.Println(colorstring.Green("That's all."))
	fmt.Println("You just have to upload the found code signing files and you'll be good to go!")
}

func initExportOutputDir() (string, error) {
	absExportOutputDirPath, err := pathutil.AbsPath(confExportOutputDirPath)
	log.Debugf("absExportOutputDirPath: %s", absExportOutputDirPath)
	if err != nil {
		return absExportOutputDirPath, fmt.Errorf("Failed to determin Absolute path of export dir: %s", confExportOutputDirPath)
	}
	if exist, err := pathutil.IsDirExists(absExportOutputDirPath); err != nil {
		return absExportOutputDirPath, fmt.Errorf("Failed to determin whether the export directory already exists: %s", err)
	} else if !exist {
		if err := os.Mkdir(absExportOutputDirPath, 0777); err != nil {
			return absExportOutputDirPath, fmt.Errorf("Failed to create export output directory at path: %s | error: %s", absExportOutputDirPath, err)
		}
	} else {
		log.Infof("Export output dir already exists at path: %s", absExportOutputDirPath)
	}
	return absExportOutputDirPath, nil
}
