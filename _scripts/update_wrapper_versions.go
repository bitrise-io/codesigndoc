package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/bitrise-io/codesigndoc/version"
	"github.com/bitrise-io/go-utils/log"
)

// VersionInfo ...
type VersionInfo struct {
	ScanCmd string
	Version string
}

func substituteVersionInfo(tmpl *template.Template, data VersionInfo, targetPth string) error {
	out, err := os.Create(targetPth)
	if err != nil {
		return fmt.Errorf("failed to open file for write, error: %s", err)
	}
	if err = tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("%s", err)
	}
	if err := out.Close(); err != nil {
		return err
	}
	return nil
}

func main() {
	log.Infof("Only update wrapper versions when release is availabe.")

	tmpl := template.Must(template.ParseFiles("install_wrap.sh.template"))

	if err := substituteVersionInfo(tmpl,
		VersionInfo{
			ScanCmd: "xcode",
			Version: version.VERSION,
		},
		"install_wrap-xcode.sh",
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	// for compatibility
	if err := substituteVersionInfo(tmpl,
		VersionInfo{
			ScanCmd: "xcode",
			Version: version.VERSION,
		},
		"install_wrap.sh",
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	if err := substituteVersionInfo(tmpl,
		VersionInfo{
			ScanCmd: "xcodeuitests",
			Version: version.VERSION,
		},
		"install_wrap-xcode-uitests.sh",
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
