## 0.9.7 (not yet released)

- FIX : Typo fix
- LOG : Separator style fix


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
