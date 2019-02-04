# Bitrise Init Tool

Initialize bitrise config, step template or plugin template

## How to build this project
Project is written in [Go](https://golang.org/) language and
uses [dep](github.com/golang/dep/cmd/dep) as dependency management tool.

You can build this project using sequence of `go` commands or refer to [bitrise.yml](./bitrise.yml) file,
which contains workflows for this project.

You can run `bitrise` workflows on your local machine using [bitrise CLI](https://www.bitrise.io/cli).

Before you start, make sure
- `$HOME/go/bin` (or `$GOPATH/bin` in case of custom go workspace) is added to `$PATH`
- `Ruby >= 2.2.2` version is installed
- `bundler` gem installed

**How to build the project using bitrise workflows**

Please check available workflows in [bitrise.yml](./bitrise.yml).
`bitrise --ci run ci` will execute `ci` workflow which consists of `prepare/build/run tests` stages.

**How to build the project using Go commands**
- `go build` command builds the project and generates `bitrise-init` binary at `$HOME/go/bin/bitrise-init`  (or `$GOPATH/bin/bitrise-init` in case of custom go workspace).
- `go test ./...` command runs unit tests in every project folder/subfolder.
- `go test -v ./_tests/integration/...` command runs integration tests. This command requires `INTEGRATION_TEST_BINARY_PATH=$HOME/go/bin/bitrise-init` (or `INTEGRATION_TEST_BINARY_PATH=$GOPATH/bin/bitrise-init` in case of custom go workspace) environment variable.

## How to release new bitrise-init version

- update the step versions in steps/const.go
- bump `version` in version/version.go
- commit these changes & open PR
- merge to master
- create tag with the new version
- test the generated release and its binaries

__Update manual config on website__

- use the generated binaries in `./_bin/` directory to generate the manual config by calling: `BIN_PATH --ci manual-config` this will generate the manual.config.yml at: `CURRENT_DIR/_defaults/result.yml`
- throw the generated `result.yml` to the frontend team, to update the manual-config on the website
- once they put the new config in the website project, check the git changes to make sure, everything looks great

__Update the [project-scanner step](https://github.com/bitrise-steplib/steps-project-scanner)__

- update bitrise-init dependency
- share a new version into the steplib (check the [README.md](https://github.com/bitrise-steplib/steps-project-scanner/blob/master/README.md))

__Update the [bitrise init plugin]((https://github.com/bitrise-core/bitrise-plugins-init))__

- update bitrise-init dependency
- release a new version (check the [README.md](https://github.com/bitrise-core/bitrise-plugins-init/blob/master/README.md))