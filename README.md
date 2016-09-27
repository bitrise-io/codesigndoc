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
1. Runs a clean Xcode Archive on your Xcode project.
1. From the Xcode logs it collects the Code Signing settings Xcode used during the Archive.
1. Prints the list of required code signing files.
1. Prints the Team ID and App/Bundle ID info of the project.
1. Optionally it can also search for, and export these files.


## Install / Run

### One-liner

Just open up your `Terminal.app` on OS X, copy-paste this into it and
hit Enter to run:

For `Xcode` project (project or workspace):

```
bash -l -c "$(curl -sfL https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap-xcode.sh)"
```

For `Xamarin` project (solution):

```
bash -l -c "$(curl -sfL https://raw.githubusercontent.com/bitrise-tools/codesigndoc/master/_scripts/install_wrap-xamarin.sh)"
```


### Manual install & run

1. download the current release - it's a single, stand-alone binary
    * example (__don't forget to replace the `VERSIONNUMBER` in the URL!__): `curl -sfL https://github.com/bitrise-tools/codesigndoc/releases/download/VERSIONNUMBER/codesigndoc-Darwin-x86_64 > ./codesigndoc`
2. `chmod +x` it, so you can run it
    * if you followed the previous example: `chmod +x ./codesigndoc`
3. run the `scan` command of the tool
    * if you followed the previous examples:
        * Xcode project scanner: `./codesigndoc scan xcode`
        * Xamarin project scanner: `./codesigndoc scan xamarin`
