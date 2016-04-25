# codesigndoc

Your friendly iOS Code Signing Doctor.

Using this tool is as easy as running `codesigndoc scan` and following the guide
it prints. At the end of the process you'll have all the code signing files
(`.p12` Identity file including the Certificate and Private Key, and the required Provisioning Profiles)
required to do a successful Xcode Archive of your Xcode project.

This tool can also help detecting code signing issues with your project,
for example it prints a warning if your build requires multiple Code Signing Identities
in order to complete an Xcode Archive.

What this tool does:

1. Gathers all information required to do a clean Xcode Archive of your project.
2. Runs a clean Xcode Archive on your Xcode project.
3. From the Xcode logs it collects the Code Signing settings Xcode used during the Archive.
4. Prints the list of required code signing files.
5. Optionally it can also search for, and export these files.


## Install / Run

### One-liner

Just open up your `Terminal.app` on OS X, copy-paste this into it and
hit Enter to run:

```
bash -l -c "$(curl -sfL https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap.sh)"
```


### Manual install & run

You can follow the steps in the [https://github.com/bitrise-tools/codesigndoc/blob/master/_scripts/install_wrap.sh](https://github.com/bitrise-tools/codesigndoc/blob/master/_scripts/install_wrap.sh)
install script file.

In short:

1. download the current release - it's a single, stand-alone binary
  * example: `curl -sfL https://github.com/bitrise-tools/codesigndoc/releases/download/0.9.2/codesigndoc-Darwin-x86_64 > ./codesigndoc`
  * make sure that you get the URL of the latest release - just replace version number (`0.9.2` in this example) in the URL with the latest release's version number
2. `chmod +x` it, so you can run it
  * if you followed the previous example: `chmod +x ./codesigndoc`
3. run the `scan` command of the tool
  * if you followed the previous examples: `./codesigndoc scan`


## TODO

- List files by project, which project used what
- Check if Certificate is revoked
