
### 0.9.8 - 0.9.7 (2016 Apr 25)

* [04ef48b] Viktor Benei - install wrap script fix (2016 Apr 25)
* [0cdf10e] Viktor Benei - v0.9.8, with updated install wrap script (2016 Apr 25)
* [1fa67db] Viktor Benei - save xcodebuild output into a debug log file (2016 Apr 25)
* [3e820d8] Viktor Benei - better install_wrap.sh script (template) (2016 Apr 25)
* [1001373] Viktor Benei - always print the xcodebuild command (2016 Apr 25)
* [4272309] Viktor Benei - generic update for bitrise.yml and .gitignore (2016 Apr 25)
* [5f94fe1] Viktor Benei - godeps update - switched to Go 1.6 vendor (2016 Apr 25)
* [97610ec] Viktor Benei - TODO (2016 Mar 09)


### 0.9.7 (2016 Feb 15)

* removed "create-test-binaries" workflow
* FIX : typo
* LOG : log separator style fix
* changelog
* use the new AskForPath goinp function for getting the project/workspace path, instead of the generic "string" version
* Godeps update
* releaseman config


## 0.9.6

- LOG : Log color revisions & additional highlights
- LOG : Better highlight for more important Warning lines


## 0.9.5

- FIX : Identities enumeration fix
- FIX : typos
- FIX : Certificate label is force converted to UTF8 before using it, to not to break in case there's a non UTF8 character
- Don't fail, just print warning if multiple Identities found for a single search (e.g. in case you have previous, revoked versions of the Certificate in your Keychain)
- Certificate filtering: export only non-expired, valid by date Certificates


## 0.9.4

- progress indicator: print a `.` every second, to indicate the Archive is still running


## 0.9.3

- log: newline fix


## 0.9.2

- `--allow-export` flag : automatically allow exporting of discovered files.
- fixed format of exported `.p12`, to be the same as if you manually export it from `Keychain Access.app`


## 0.9.1

- First public (beta) version

