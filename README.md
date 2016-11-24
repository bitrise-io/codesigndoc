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


## Manually finding the required base code signing files for an Xcode project or workspace

If you'd want to manually check which files are __required__ for archiving your project (regardless of the distribution type!), you have to do a clean archive **on your Mac**, using Xcode's command line tool (`xcodebuild`) and check the logs.
The easiest way is open the Terminal app, `cd` into the directory where your Xcode project/workspace file is located, and do a clean archive from Terminal.

Performing a clean archive from Terminal is as easy as running this command (*on your Mac*) if you use an Xcode Workspace: `xcodebuild -workspace "YOUR.xcworkspace" -scheme "a Shared scheme" clean archive` or this one if you use an Xcode Project: `xcodebuild -project "YOUR.xcodeproj" -scheme "a Shared scheme" clean archive`

In the output you'll see code signing infos, namely you should search for the text `Signing Identity` which is followed by a `Provisioning Profile` line. There might be more than one configuration in the log - these are the configurations used by Xcode on your Mac when you do an Archive.

To make an Xcode Archive work _on any Mac_, you need the same Code Signing Identity (certificate) and Provisioning Profile(s). T_he signing Identities (certificates) and Provisioning Profiles present in the log are __required__, regardless of the final distribution type you use._

To run the `xcodebuild` command and only show these lines you can add the postfix: ` | grep -i -e 'Signing Identity' -e 'Provisioning Profile' -e '` to your call, for example: `xcodebuild -workspace "YOUR.xcworkspace" -scheme "a Shared scheme" clean archive | grep -i -e 'Signing Identity' -e 'Provisioning Profile' -e '`. This will run the exact same command, but will filter out every other text in the output except these lines you're searching for.

By running this command you'll see an output similar to:

```
Signing Identity:     "iPhone Developer: Viktor Benei (F...7)"
Provisioning Profile: "BuildAnything"
                      (9.......-....-....-....-...........0)
```

If you see more than one `Signing Identity` or `Provisioning Profile` line that means that Xcode had to switch between code signing configurations to be able to create your archive. __All of the listed certificates & provisioning profiles have to be available to create an archive of your project__ with your current code signing settings.
