package project

const iosTestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">iPhoneSimulator</Platform>
    <ProjectTypeGuids>{FEACFBD2-3405-455C-9665-78FE426C6842};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}</ProjectTypeGuids>
    <ProjectGuid>{90F3C584-FD69-4926-9903-6B9771847782}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>CreditCardValidator.iOS</RootNamespace>
    <IPhoneResourcePrefix>Resources</IPhoneResourcePrefix>
    <AssemblyName>CreditCardValidator.iOS</AssemblyName>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|iPhoneSimulator' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\iPhoneSimulator\Debug</OutputPath>
    <DefineConstants>DEBUG;ENABLE_TEST_CLOUD;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <MtouchArch>i386</MtouchArch>
    <MtouchLink>None</MtouchLink>
    <MtouchDebug>true</MtouchDebug>
    <MtouchProfiling>true</MtouchProfiling>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|iPhone' ">
    <Optimize>true</Optimize>
    <OutputPath>bin\iPhone\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <CodesignEntitlements>Entitlements.plist</CodesignEntitlements>
    <MtouchArch>ARMv7, ARM64</MtouchArch>
    <ConsolePause>false</ConsolePause>
    <CodesignKey>iPhone Developer: Bitrise Bot (VV2J4SV8V4)</CodesignKey>
    <CodesignProvision>225561e6-3526-4edc-a046-7e0fa49eb4fe</CodesignProvision>
    <IpaPackageName>
    </IpaPackageName>
    <BuildIpa>true</BuildIpa>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|iPhoneSimulator' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\iPhoneSimulator\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <MtouchArch>i386</MtouchArch>
    <ConsolePause>false</ConsolePause>
    <MtouchLink>None</MtouchLink>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|iPhone' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\iPhone\Debug</OutputPath>
    <DefineConstants>DEBUG;ENABLE_TEST_CLOUD;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <MtouchArch>ARMv7, ARM64</MtouchArch>
    <CodesignEntitlements>Entitlements.plist</CodesignEntitlements>
    <MtouchProfiling>true</MtouchProfiling>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchDebug>true</MtouchDebug>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchI18n>
    </MtouchI18n>
    <IpaPackageName>
    </IpaPackageName>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="System.Xml" />
    <Reference Include="System.Core" />
    <Reference Include="Xamarin.iOS" />
    <Reference Include="Calabash">
      <HintPath>..\packages\Xamarin.TestCloud.Agent.0.19.1\lib\Xamarin.iOS10\Calabash.dll</HintPath>
    </Reference>
  </ItemGroup>
  <ItemGroup>
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Contents.json" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-60%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-Small.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-Small%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-Spotlight-40.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-Spotlight-40%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-76.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-76%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-72.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-72%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-Small-50.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Icon-Small-50%402x.png" />
  </ItemGroup>
  <ItemGroup>
    <InterfaceDefinition Include="Resources\LaunchScreen.xib" />
    <InterfaceDefinition Include="Main.storyboard" />
  </ItemGroup>
  <ItemGroup>
    <None Include="Info.plist" />
    <None Include="Entitlements.plist" />
    <None Include="packages.config" />
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Main.cs" />
    <Compile Include="AppDelegate.cs" />
    <Compile Include="ViewController.cs" />
    <Compile Include="ViewController.designer.cs">
      <DependentUpon>ViewController.cs</DependentUpon>
    </Compile>
    <Compile Include="ValidCreditCardController.cs" />
    <Compile Include="ValidCreditCardController.designer.cs">
      <DependentUpon>ValidCreditCardController.cs</DependentUpon>
    </Compile>
  </ItemGroup>
  <Import Project="$(MSBuildExtensionsPath)\Xamarin\iOS\Xamarin.iOS.CSharp.targets" />
  <ItemGroup>
    <ProjectReference Include="..\CreditCardValidator\CreditCardValidator.csproj">
      <Project>{99A825A6-6F99-4B94-9F65-E908A6347F1E}</Project>
      <Name>CreditCardValidator</Name>
    </ProjectReference>
  </ItemGroup>
</Project>`

const androidTestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectTypeGuids>{EFBA0AD7-5A72-4C68-AF49-83D382785DCF};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}</ProjectTypeGuids>
    <ProjectGuid>{9D1D32A3-D13F-4F23-B7D4-EF9D52B06E60}</ProjectGuid>
    <OutputType>Library</OutputType>
    <RootNamespace>CreditCardValidator.Droid</RootNamespace>
    <MonoAndroidAssetsPrefix>Assets</MonoAndroidAssetsPrefix>
    <MonoAndroidResourcePrefix>Resources</MonoAndroidResourcePrefix>
    <AndroidResgenClass>Resource</AndroidResgenClass>
    <AndroidResgenFile>Resources\Resource.designer.cs</AndroidResgenFile>
    <AndroidApplication>True</AndroidApplication>
    <AndroidUseLatestPlatformSdk>False</AndroidUseLatestPlatformSdk>
    <AssemblyName>CreditCardValidator.Droid</AssemblyName>
    <AndroidManifest>Properties\AndroidManifest.xml</AndroidManifest>
    <TargetFrameworkVersion>v4.4</TargetFrameworkVersion>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|AnyCPU' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <AndroidLinkMode>None</AndroidLinkMode>
    <ConsolePause>false</ConsolePause>
    <AndroidSigningKeyStore>bitrise-android-keystore.jks</AndroidSigningKeyStore>
    <AndroidSigningStorePass>bitrise</AndroidSigningStorePass>
    <AndroidSigningKeyAlias>bitrise-alias</AndroidSigningKeyAlias>
    <AndroidSigningKeyPass>bitrise</AndroidSigningKeyPass>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <AndroidSupportedAbis>armeabi-v7a;x86</AndroidSupportedAbis>
    <AndroidSigningKeyStore>bitrise-android-keystore.jks</AndroidSigningKeyStore>
    <AndroidSigningStorePass>bitrise</AndroidSigningStorePass>
    <AndroidSigningKeyAlias>bitrise-alias</AndroidSigningKeyAlias>
    <AndroidSigningKeyPass>bitrise</AndroidSigningKeyPass>
    <EmbedAssembliesIntoApk>True</EmbedAssembliesIntoApk>
    <AndroidUseSharedRuntime>false</AndroidUseSharedRuntime>
    <AndroidKeyStore>True</AndroidKeyStore>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="System.Xml" />
    <Reference Include="System.Core" />
    <Reference Include="Mono.Android" />
  </ItemGroup>
  <ItemGroup>
    <Compile Include="MainActivity.cs" />
    <Compile Include="Resources\Resource.designer.cs" />
    <Compile Include="Properties\AssemblyInfo.cs" />
    <Compile Include="CreditCardValidationSuccess.cs" />
  </ItemGroup>
  <ItemGroup>
    <None Include="Resources\AboutResources.txt" />
    <None Include="Properties\AndroidManifest.xml" />
    <None Include="Assets\AboutAssets.txt" />
  </ItemGroup>
  <ItemGroup>
    <AndroidResource Include="Resources\layout\Main.axml" />
    <AndroidResource Include="Resources\values\Strings.xml" />
    <AndroidResource Include="Resources\drawable-hdpi\Icon.png" />
    <AndroidResource Include="Resources\drawable-mdpi\Icon.png" />
    <AndroidResource Include="Resources\drawable-xhdpi\Icon.png" />
    <AndroidResource Include="Resources\drawable-xxhdpi\Icon.png" />
    <AndroidResource Include="Resources\drawable-xxxhdpi\Icon.png" />
    <AndroidResource Include="Resources\layout\CreditCardValidationSuccess.axml" />
  </ItemGroup>
  <ItemGroup>
    <ProjectReference Include="..\CreditCardValidator\CreditCardValidator.csproj">
      <Project>{99A825A6-6F99-4B94-9F65-E908A6347F1E}</Project>
      <Name>CreditCardValidator</Name>
    </ProjectReference>
  </ItemGroup>
  <Import Project="$(MSBuildExtensionsPath)\Xamarin\Android\Xamarin.Android.CSharp.targets" />
</Project>`

const macTestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectTypeGuids>{A3F8F2AB-B479-4A4A-A458-A89E7DC349F1};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}</ProjectTypeGuids>
    <ProjectGuid>{4DA5EAC6-6F80-4FEC-AF81-194210F10B51}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>Hello_Mac</RootNamespace>
    <MonoMacResourcePrefix>Resources</MonoMacResourcePrefix>
    <AssemblyName>Hello_Mac</AssemblyName>
    <TargetFrameworkVersion>v2.0</TargetFrameworkVersion>
    <TargetFrameworkIdentifier>Xamarin.Mac</TargetFrameworkIdentifier>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|AnyCPU' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <Profiling>true</Profiling>
    <UseRefCounting>true</UseRefCounting>
    <UseSGen>true</UseSGen>
    <IncludeMonoRuntime>false</IncludeMonoRuntime>
    <CreatePackage>false</CreatePackage>
    <CodeSigningKey>Mac Developer</CodeSigningKey>
    <EnableCodeSigning>false</EnableCodeSigning>
    <EnablePackageSigning>false</EnablePackageSigning>
    <PackageSigningKey>Developer ID Installer</PackageSigningKey>
    <CodeSignEntitlements>Entitlements.plist</CodeSignEntitlements>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <LinkMode>SdkOnly</LinkMode>
    <Profiling>false</Profiling>
    <UseRefCounting>true</UseRefCounting>
    <UseSGen>true</UseSGen>
    <IncludeMonoRuntime>true</IncludeMonoRuntime>
    <CreatePackage>true</CreatePackage>
    <CodeSigningKey>Mac Developer</CodeSigningKey>
    <EnableCodeSigning>true</EnableCodeSigning>
    <EnablePackageSigning>false</EnablePackageSigning>
    <PackageSigningKey>Developer ID Installer</PackageSigningKey>
    <XamMacArch>x86_64</XamMacArch>
    <CodeSignProvision>Automatic</CodeSignProvision>
    <CodeSignEntitlements>Entitlements.plist</CodeSignEntitlements>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="System.Core" />
    <Reference Include="Xamarin.Mac" />
  </ItemGroup>
  <ItemGroup>
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\Contents.json" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-128.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-128%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-16.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-16%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-256.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-256%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-32.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-32%402x.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-512.png" />
    <ImageAsset Include="Resources\Images.xcassets\AppIcons.appiconset\AppIcon-512%402x.png" />
  </ItemGroup>
  <ItemGroup>
    <None Include="Info.plist" />
    <None Include="Entitlements.plist" />
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Main.cs" />
    <Compile Include="AppDelegate.cs" />
    <Compile Include="ViewController.cs" />
    <Compile Include="ViewController.designer.cs">
      <DependentUpon>ViewController.cs</DependentUpon>
    </Compile>
  </ItemGroup>
  <ItemGroup>
    <InterfaceDefinition Include="Main.storyboard" />
  </ItemGroup>
  <Import Project="$(MSBuildExtensionsPath)\Xamarin\Mac\Xamarin.Mac.CSharp.targets" />
</Project>`

const tvTestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">iPhoneSimulator</Platform>
    <ProjectGuid>{51D9C362-2997-4029-B38F-06C36F17056E}</ProjectGuid>
    <ProjectTypeGuids>{06FA79CB-D6CD-4721-BB4B-1BD202089C55};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}</ProjectTypeGuids>
    <OutputType>Exe</OutputType>
    <RootNamespace>tvos</RootNamespace>
    <AssemblyName>tvos</AssemblyName>
    <IPhoneResourcePrefix>Resources</IPhoneResourcePrefix>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|iPhoneSimulator' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\iPhoneSimulator\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchDebug>true</MtouchDebug>
    <MtouchFastDev>true</MtouchFastDev>
    <MtouchProfiling>true</MtouchProfiling>
    <MtouchUseSGen>true</MtouchUseSGen>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchLink>None</MtouchLink>
    <MtouchArch>x86_64</MtouchArch>
    <MtouchHttpClientHandler>HttpClientHandler</MtouchHttpClientHandler>
    <MtouchTlsProvider>Default</MtouchTlsProvider>
    <PlatformTarget>x86</PlatformTarget>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|iPhone' ">
    <Optimize>true</Optimize>
    <OutputPath>bin\iPhone\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchUseLlvm>true</MtouchUseLlvm>
    <MtouchUseSGen>true</MtouchUseSGen>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchFloat32>true</MtouchFloat32>
    <MtouchEnableBitcode>true</MtouchEnableBitcode>
    <CodesignEntitlements>Entitlements.plist</CodesignEntitlements>
    <MtouchLink>SdkOnly</MtouchLink>
    <MtouchArch>ARM64</MtouchArch>
    <MtouchHttpClientHandler>HttpClientHandler</MtouchHttpClientHandler>
    <MtouchTlsProvider>Default</MtouchTlsProvider>
    <PlatformTarget>x86</PlatformTarget>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|iPhoneSimulator' ">
    <Optimize>true</Optimize>
    <OutputPath>bin\iPhoneSimulator\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchUseSGen>true</MtouchUseSGen>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchLink>None</MtouchLink>
    <MtouchArch>x86_64</MtouchArch>
    <MtouchHttpClientHandler>HttpClientHandler</MtouchHttpClientHandler>
    <MtouchTlsProvider>Default</MtouchTlsProvider>
    <PlatformTarget>x86</PlatformTarget>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|iPhone' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\iPhone\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <CodesignKey>iPhone Developer</CodesignKey>
    <DeviceSpecificBuild>true</DeviceSpecificBuild>
    <MtouchDebug>true</MtouchDebug>
    <MtouchFastDev>true</MtouchFastDev>
    <MtouchProfiling>true</MtouchProfiling>
    <MtouchUseSGen>true</MtouchUseSGen>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchFloat32>true</MtouchFloat32>
    <CodesignEntitlements>Entitlements.plist</CodesignEntitlements>
    <MtouchLink>SdkOnly</MtouchLink>
    <MtouchArch>ARM64</MtouchArch>
    <MtouchHttpClientHandler>HttpClientHandler</MtouchHttpClientHandler>
    <MtouchTlsProvider>Default</MtouchTlsProvider>
    <PlatformTarget>x86</PlatformTarget>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="System.Xml" />
    <Reference Include="System.Core" />
    <Reference Include="Xamarin.TVOS" />
  </ItemGroup>
  <ItemGroup>
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Back.imagestacklayer\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Back.imagestacklayer\Content.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Front.imagestacklayer\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Front.imagestacklayer\Content.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Middle.imagestacklayer\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Large.imagestack\Middle.imagestacklayer\Content.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Back.imagestacklayer\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Back.imagestacklayer\Content.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Front.imagestacklayer\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Front.imagestacklayer\Content.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Middle.imagestacklayer\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\App Icon - Small.imagestack\Middle.imagestacklayer\Content.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\App Icon &amp; Top Shelf Image.brandassets\Top Shelf Image.imageset\Contents.json" />
    <ImageAsset Include="Assets.xcassets\LaunchImages.launchimage\Contents.json" />
    <ImageAsset Include="Assets.xcassets\Contents.json" />
  </ItemGroup>
  <ItemGroup>
    <Folder Include="Resources\" />
  </ItemGroup>
  <ItemGroup>
    <None Include="Info.plist" />
    <None Include="Entitlements.plist" />
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Main.cs" />
    <Compile Include="AppDelegate.cs" />
    <Compile Include="ViewController.cs" />
    <Compile Include="ViewController.designer.cs">
      <DependentUpon>ViewController.cs</DependentUpon>
    </Compile>
  </ItemGroup>
  <ItemGroup>
    <InterfaceDefinition Include="Main.storyboard" />
  </ItemGroup>
  <Import Project="$(MSBuildExtensionsPath)\Xamarin\TVOS\Xamarin.TVOS.CSharp.targets" />
</Project>`

const xamarinUITestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectGuid>{BA48743D-06F3-4D2D-ACFD-EE2642CE155A}</ProjectGuid>
    <OutputType>Library</OutputType>
    <RootNamespace>CreditCardValidator.iOS.UITests</RootNamespace>
    <AssemblyName>CreditCardValidator.iOS.UITests</AssemblyName>
    <TargetFrameworkVersion>v4.5</TargetFrameworkVersion>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|AnyCPU' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="nunit.framework">
      <HintPath>..\packages\NUnit.2.6.4\lib\nunit.framework.dll</HintPath>
    </Reference>
    <Reference Include="Xamarin.UITest">
      <HintPath>..\packages\Xamarin.UITest.1.3.8\lib\Xamarin.UITest.dll</HintPath>
    </Reference>
  </ItemGroup>
  <ItemGroup>
    <ProjectReference Include="..\CreditCardValidator.iOS\CreditCardValidator.iOS.csproj">
      <Project>{90F3C584-FD69-4926-9903-6B9771847782}</Project>
      <Name>CreditCardValidator.iOS</Name>
      <ReferenceOutputAssembly>False</ReferenceOutputAssembly>
      <Private>False</Private>
    </ProjectReference>
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Tests.cs" />
  </ItemGroup>
  <Import Project="$(MSBuildBinPath)\Microsoft.CSharp.targets" />
  <ItemGroup>
    <None Include="packages.config" />
  </ItemGroup>
</Project>`

const nunitTestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectGuid>{ED150913-76EB-446F-8B78-DC77E5795703}</ProjectGuid>
    <OutputType>Library</OutputType>
    <RootNamespace>CreditCardValidator.iOS.NunitTests</RootNamespace>
    <AssemblyName>CreditCardValidator.iOS.NunitTests</AssemblyName>
    <TargetFrameworkVersion>v4.5</TargetFrameworkVersion>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|AnyCPU' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="nunit.framework">
      <HintPath>..\packages\NUnit.2.6.4\lib\nunit.framework.dll</HintPath>
    </Reference>
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Test.cs" />
  </ItemGroup>
  <Import Project="$(MSBuildBinPath)\Microsoft.CSharp.targets" />
  <ItemGroup>
    <None Include="packages.config" />
  </ItemGroup>
</Project>`

const nunitLiteTestProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">iPhoneSimulator</Platform>
    <ProjectTypeGuids>{FEACFBD2-3405-455C-9665-78FE426C6842};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}</ProjectTypeGuids>
    <ProjectGuid>{95615CA5-0D75-4389-A6E0-78309A686712}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>CreditCardValidator.iOS.NunitLiteTests</RootNamespace>
    <IPhoneResourcePrefix>Resources</IPhoneResourcePrefix>
    <AssemblyName>CreditCardValidator.iOS.NunitLiteTests</AssemblyName>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|iPhoneSimulator' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\iPhoneSimulator\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <MtouchArch>i386</MtouchArch>
    <MtouchLink>None</MtouchLink>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchUseSGen>true</MtouchUseSGen>
    <MtouchFastDev>true</MtouchFastDev>
    <MtouchDebug>true</MtouchDebug>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchProfiling>true</MtouchProfiling>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|iPhone' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\iPhone\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <MtouchArch>ARMv7, ARM64</MtouchArch>
    <CodesignEntitlements>Entitlements.plist</CodesignEntitlements>
    <MtouchFloat32>true</MtouchFloat32>
    <MtouchUseSGen>true</MtouchUseSGen>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|iPhoneSimulator' ">
    <DebugType>full</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\iPhoneSimulator\Release</OutputPath>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <MtouchArch>i386</MtouchArch>
    <MtouchLink>None</MtouchLink>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <CodesignKey>iPhone Developer</CodesignKey>
    <MtouchUseSGen>true</MtouchUseSGen>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|iPhone' ">
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\iPhone\Debug</OutputPath>
    <DefineConstants>DEBUG;</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <ConsolePause>false</ConsolePause>
    <MtouchArch>ARMv7, ARM64</MtouchArch>
    <CodesignEntitlements>Entitlements.plist</CodesignEntitlements>
    <MtouchFloat32>true</MtouchFloat32>
    <CodesignKey>iPhone Developer</CodesignKey>
    <DeviceSpecificBuild>true</DeviceSpecificBuild>
    <MtouchDebug>true</MtouchDebug>
    <MtouchUseSGen>true</MtouchUseSGen>
    <MtouchUseRefCounting>true</MtouchUseRefCounting>
    <MtouchProfiling>true</MtouchProfiling>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="System.Xml" />
    <Reference Include="System.Core" />
    <Reference Include="Xamarin.iOS" />
    <Reference Include="MonoTouch.NUnitLite" />
  </ItemGroup>
  <ItemGroup>
    <Folder Include="Resources\" />
  </ItemGroup>
  <ItemGroup>
    <None Include="Info.plist" />
    <None Include="Entitlements.plist" />
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Main.cs" />
    <Compile Include="UnitTestAppDelegate.cs" />
    <Compile Include="TestsSample.cs" />
  </ItemGroup>
  <Import Project="$(MSBuildExtensionsPath)\Xamarin\iOS\Xamarin.iOS.CSharp.targets" />
</Project>`
