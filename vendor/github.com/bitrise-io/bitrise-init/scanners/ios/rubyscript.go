package ios

import (
	"errors"
	"path"

	"os"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
)

func runRubyScriptForOutput(scriptContent, gemfileContent string, withEnvs []string) (string, error) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__bitrise-init__")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.TErrorf("Failed to remove tmp dir (%s), error: %s", tmpDir, err)
		}
	}()

	// Write Gemfile to file and install
	if gemfileContent != "" {
		gemfilePth := path.Join(tmpDir, "Gemfile")
		if err := fileutil.WriteStringToFile(gemfilePth, gemfileContent); err != nil {
			return "", err
		}

		cmd := command.New("bundle", "install")
		cmd.SetDir(tmpDir)

		withEnvs = append(withEnvs, "BUNDLE_GEMFILE="+gemfilePth)
		cmd.AppendEnvs(withEnvs...)

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			if errorutil.IsExitStatusError(err) {
				return "", errors.New(out)
			}
			return "", err
		}
	}

	// Write script to file and run
	rubyScriptPth := path.Join(tmpDir, "script.rb")
	if err := fileutil.WriteStringToFile(rubyScriptPth, scriptContent); err != nil {
		return "", err
	}

	var cmd *command.Model

	if gemfileContent != "" {
		cmd = command.New("bundle", "exec", "ruby", rubyScriptPth)
	} else {
		cmd = command.New("ruby", rubyScriptPth)
	}

	// Set the temp dir as working dir, so the project defined `.ruby-version` does not cause ruby resolution to fail:
	// [ ... ] ruby script failed, error: rbenv: version `2.7.4' is not installed (set by /[ ... ]/MyTestApp/.ruby-version)
	cmd.SetDir(tmpDir)

	if len(withEnvs) > 0 {
		cmd.AppendEnvs(withEnvs...)
	}

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		if errorutil.IsExitStatusError(err) {
			return "", errors.New(out)
		}
		return "", err
	}

	return out, nil
}
