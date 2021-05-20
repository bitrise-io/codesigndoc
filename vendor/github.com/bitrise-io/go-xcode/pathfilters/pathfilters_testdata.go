package pathfilters

const testMacOSPbxprojContent = `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 46;
	objects = {

/* Begin PBXBuildFile section */
		1302F7441D95BA4A005CE678 /* Session.swift in Sources */ = {isa = PBXBuildFile; fileRef = 1302F7431D95BA4A005CE678 /* Session.swift */; };
		130E6BBE1D95BBB4009D3C78 /* Command.swift in Sources */ = {isa = PBXBuildFile; fileRef = 138F9EE61D8E7ABC00515FCA /* Command.swift */; };
		131A3D9A1D90543F002DAF99 /* Realm.framework in Frameworks */ = {isa = PBXBuildFile; fileRef = 131A3D961D9053BB002DAF99 /* Realm.framework */; };
		131A3D9B1D90543F002DAF99 /* Realm.framework in Embed Frameworks */ = {isa = PBXBuildFile; fileRef = 131A3D961D9053BB002DAF99 /* Realm.framework */; settings = {ATTRIBUTES = (CodeSignOnCopy, RemoveHeadersOnCopy, ); }; };
		131A3D9C1D90543F002DAF99 /* RealmSwift.framework in Frameworks */ = {isa = PBXBuildFile; fileRef = 131A3D971D9053BB002DAF99 /* RealmSwift.framework */; };
		131A3D9D1D90543F002DAF99 /* RealmSwift.framework in Embed Frameworks */ = {isa = PBXBuildFile; fileRef = 131A3D971D9053BB002DAF99 /* RealmSwift.framework */; settings = {ATTRIBUTES = (CodeSignOnCopy, RemoveHeadersOnCopy, ); }; };
		131A3DA61D9060AA002DAF99 /* Yaml.framework in Embed Frameworks */ = {isa = PBXBuildFile; fileRef = 131A3DA11D90609E002DAF99 /* Yaml.framework */; settings = {ATTRIBUTES = (CodeSignOnCopy, RemoveHeadersOnCopy, ); }; };
		131A3DA71D9060C0002DAF99 /* Yaml.framework in Frameworks */ = {isa = PBXBuildFile; fileRef = 131A3DA11D90609E002DAF99 /* Yaml.framework */; };
		131A3DAB1D906BFA002DAF99 /* BitriseTool.swift in Sources */ = {isa = PBXBuildFile; fileRef = 131A3DAA1D906BFA002DAF99 /* BitriseTool.swift */; };
		131A3DAD1D906D54002DAF99 /* Bitrise.swift in Sources */ = {isa = PBXBuildFile; fileRef = 131A3DAC1D906D54002DAF99 /* Bitrise.swift */; };
		131A3DAF1D906DAB002DAF99 /* Envman.swift in Sources */ = {isa = PBXBuildFile; fileRef = 131A3DAE1D906DAB002DAF99 /* Envman.swift */; };
		131A3DB11D906DBF002DAF99 /* Stepman.swift in Sources */ = {isa = PBXBuildFile; fileRef = 131A3DB01D906DBF002DAF99 /* Stepman.swift */; };
		131ACE731D93054B007E71E9 /* ToolManager.swift in Sources */ = {isa = PBXBuildFile; fileRef = 131ACE721D93054B007E71E9 /* ToolManager.swift */; };
		131ACE761D9323E3007E71E9 /* Version.swift in Sources */ = {isa = PBXBuildFile; fileRef = 131ACE751D9323E3007E71E9 /* Version.swift */; };
		13A95E821D966F040061B54F /* BashSessionViewController.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13A95E811D966F040061B54F /* BashSessionViewController.swift */; };
		13AAE1BD1D8426FF00AEE66D /* FileManager.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13AAE1BC1D8426FF00AEE66D /* FileManager.swift */; };
		13B0B8F31D872E93006EA29C /* RealmManager.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13B0B8F21D872E93006EA29C /* RealmManager.swift */; };
		13B958CD1D89E87600D3310D /* SystemInfoViewController.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13B958CC1D89E87600D3310D /* SystemInfoViewController.swift */; };
		13B958D01D89EC7800D3310D /* RunViewController.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13B958CE1D89EC7800D3310D /* RunViewController.swift */; };
		13C989811D8319600028BA2C /* AppDelegate.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C989801D8319600028BA2C /* AppDelegate.swift */; };
		13C989851D8319600028BA2C /* Assets.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 13C989841D8319600028BA2C /* Assets.xcassets */; };
		13C989881D8319600028BA2C /* Main.storyboard in Resources */ = {isa = PBXBuildFile; fileRef = 13C989861D8319600028BA2C /* Main.storyboard */; };
		13C989931D8319600028BA2C /* BitriseStudioTests.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C989921D8319600028BA2C /* BitriseStudioTests.swift */; };
		13C9899E1D8319600028BA2C /* BitriseStudioUITests.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C9899D1D8319600028BA2C /* BitriseStudioUITests.swift */; };
		13CA73D61D84A74800B1A323 /* AddProjectViewController.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13CA73D51D84A74800B1A323 /* AddProjectViewController.swift */; };
		13CA73DE1D84B67200B1A323 /* Project.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13CA73DD1D84B67200B1A323 /* Project.swift */; };
		13E3F5531D83477300AE7C20 /* ProjectsViewController.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C989821D8319600028BA2C /* ProjectsViewController.swift */; };
		13FF5FCE1D9859EE008C7DFB /* Log.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13FF5FCD1D9859EE008C7DFB /* Log.swift */; };
		13FF5FD11D98620A008C7DFB /* String+Extensions.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13FF5FD01D98620A008C7DFB /* String+Extensions.swift */; };
		13FF5FD31D9862DA008C7DFB /* Data+Extensions.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13FF5FD21D9862DA008C7DFB /* Data+Extensions.swift */; };
		13FF5FD51D9872EC008C7DFB /* Pipe+Extensions.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13FF5FD41D9872EC008C7DFB /* Pipe+Extensions.swift */; };
/* End PBXBuildFile section */

/* Begin PBXContainerItemProxy section */
		13C9898F1D8319600028BA2C /* PBXContainerItemProxy */ = {
			isa = PBXContainerItemProxy;
			containerPortal = 13C989751D83195F0028BA2C /* Project object */;
			proxyType = 1;
			remoteGlobalIDString = 13C9897C1D83195F0028BA2C;
			remoteInfo = BitriseStudio;
		};
		13C9899A1D8319600028BA2C /* PBXContainerItemProxy */ = {
			isa = PBXContainerItemProxy;
			containerPortal = 13C989751D83195F0028BA2C /* Project object */;
			proxyType = 1;
			remoteGlobalIDString = 13C9897C1D83195F0028BA2C;
			remoteInfo = BitriseStudio;
		};
/* End PBXContainerItemProxy section */

/* Begin PBXCopyFilesBuildPhase section */
		131A3D9E1D90543F002DAF99 /* Embed Frameworks */ = {
			isa = PBXCopyFilesBuildPhase;
			buildActionMask = 2147483647;
			dstPath = "";
			dstSubfolderSpec = 10;
			files = (
				131A3D9D1D90543F002DAF99 /* RealmSwift.framework in Embed Frameworks */,
				131A3D9B1D90543F002DAF99 /* Realm.framework in Embed Frameworks */,
				131A3DA61D9060AA002DAF99 /* Yaml.framework in Embed Frameworks */,
			);
			name = "Embed Frameworks";
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXCopyFilesBuildPhase section */

/* Begin PBXFileReference section */
		1302F7431D95BA4A005CE678 /* Session.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Session.swift; sourceTree = "<group>"; };
		131A3D961D9053BB002DAF99 /* Realm.framework */ = {isa = PBXFileReference; lastKnownFileType = wrapper.framework; name = Realm.framework; path = Framworks/Realm.framework; sourceTree = "<group>"; };
		131A3D971D9053BB002DAF99 /* RealmSwift.framework */ = {isa = PBXFileReference; lastKnownFileType = wrapper.framework; name = RealmSwift.framework; path = Framworks/RealmSwift.framework; sourceTree = "<group>"; };
		131A3DA11D90609E002DAF99 /* Yaml.framework */ = {isa = PBXFileReference; lastKnownFileType = wrapper.framework; name = Yaml.framework; path = Framworks/Yaml.framework; sourceTree = "<group>"; };
		131A3DAA1D906BFA002DAF99 /* BitriseTool.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = BitriseTool.swift; sourceTree = "<group>"; };
		131A3DAC1D906D54002DAF99 /* Bitrise.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Bitrise.swift; sourceTree = "<group>"; };
		131A3DAE1D906DAB002DAF99 /* Envman.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Envman.swift; sourceTree = "<group>"; };
		131A3DB01D906DBF002DAF99 /* Stepman.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Stepman.swift; sourceTree = "<group>"; };
		131ACE721D93054B007E71E9 /* ToolManager.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = ToolManager.swift; sourceTree = "<group>"; };
		131ACE751D9323E3007E71E9 /* Version.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Version.swift; sourceTree = "<group>"; };
		138F9EE61D8E7ABC00515FCA /* Command.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Command.swift; sourceTree = "<group>"; };
		13A95E811D966F040061B54F /* BashSessionViewController.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = BashSessionViewController.swift; sourceTree = "<group>"; };
		13AAE1BC1D8426FF00AEE66D /* FileManager.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = FileManager.swift; sourceTree = "<group>"; };
		13B0B8F21D872E93006EA29C /* RealmManager.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = RealmManager.swift; sourceTree = "<group>"; };
		13B958CC1D89E87600D3310D /* SystemInfoViewController.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = SystemInfoViewController.swift; sourceTree = "<group>"; };
		13B958CE1D89EC7800D3310D /* RunViewController.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = RunViewController.swift; sourceTree = "<group>"; };
		13C9897D1D8319600028BA2C /* BitriseStudio.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = BitriseStudio.app; sourceTree = BUILT_PRODUCTS_DIR; };
		13C989801D8319600028BA2C /* AppDelegate.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = AppDelegate.swift; sourceTree = "<group>"; };
		13C989821D8319600028BA2C /* ProjectsViewController.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = ProjectsViewController.swift; sourceTree = "<group>"; };
		13C989841D8319600028BA2C /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
		13C989871D8319600028BA2C /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/Main.storyboard; sourceTree = "<group>"; };
		13C989891D8319600028BA2C /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13C9898E1D8319600028BA2C /* BitriseStudioTests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = BitriseStudioTests.xctest; sourceTree = BUILT_PRODUCTS_DIR; };
		13C989921D8319600028BA2C /* BitriseStudioTests.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = BitriseStudioTests.swift; sourceTree = "<group>"; };
		13C989941D8319600028BA2C /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13C989991D8319600028BA2C /* BitriseStudioUITests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = BitriseStudioUITests.xctest; sourceTree = BUILT_PRODUCTS_DIR; };
		13C9899D1D8319600028BA2C /* BitriseStudioUITests.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = BitriseStudioUITests.swift; sourceTree = "<group>"; };
		13C9899F1D8319600028BA2C /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13CA73D51D84A74800B1A323 /* AddProjectViewController.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = AddProjectViewController.swift; sourceTree = "<group>"; };
		13CA73DD1D84B67200B1A323 /* Project.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Project.swift; sourceTree = "<group>"; };
		13FF5FCD1D9859EE008C7DFB /* Log.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = Log.swift; sourceTree = "<group>"; };
		13FF5FD01D98620A008C7DFB /* String+Extensions.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = "String+Extensions.swift"; sourceTree = "<group>"; };
		13FF5FD21D9862DA008C7DFB /* Data+Extensions.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = "Data+Extensions.swift"; sourceTree = "<group>"; };
		13FF5FD41D9872EC008C7DFB /* Pipe+Extensions.swift */ = {isa = PBXFileReference; fileEncoding = 4; lastKnownFileType = sourcecode.swift; path = "Pipe+Extensions.swift"; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXFrameworksBuildPhase section */
		13C9897A1D83195F0028BA2C /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
				131A3D9C1D90543F002DAF99 /* RealmSwift.framework in Frameworks */,
				131A3D9A1D90543F002DAF99 /* Realm.framework in Frameworks */,
				131A3DA71D9060C0002DAF99 /* Yaml.framework in Frameworks */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C9898B1D8319600028BA2C /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C989961D8319600028BA2C /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXFrameworksBuildPhase section */

/* Begin PBXGroup section */
		1302F7421D95BA35005CE678 /* Bash */ = {
			isa = PBXGroup;
			children = (
				1302F7431D95BA4A005CE678 /* Session.swift */,
				138F9EE61D8E7ABC00515FCA /* Command.swift */,
			);
			name = Bash;
			sourceTree = "<group>";
		};
		131ACE741D9323C1007E71E9 /* Version */ = {
			isa = PBXGroup;
			children = (
				131ACE751D9323E3007E71E9 /* Version.swift */,
			);
			name = Version;
			sourceTree = "<group>";
		};
		13C989741D83195F0028BA2C = {
			isa = PBXGroup;
			children = (
				13C9897F1D8319600028BA2C /* BitriseStudio */,
				13C989911D8319600028BA2C /* BitriseStudioTests */,
				13C9899C1D8319600028BA2C /* BitriseStudioUITests */,
				13C9897E1D8319600028BA2C /* Products */,
				13CA73D71D84B5C500B1A323 /* Frameworks */,
			);
			sourceTree = "<group>";
		};
		13C9897E1D8319600028BA2C /* Products */ = {
			isa = PBXGroup;
			children = (
				13C9897D1D8319600028BA2C /* BitriseStudio.app */,
				13C9898E1D8319600028BA2C /* BitriseStudioTests.xctest */,
				13C989991D8319600028BA2C /* BitriseStudioUITests.xctest */,
			);
			name = Products;
			sourceTree = "<group>";
		};
		13C9897F1D8319600028BA2C /* BitriseStudio */ = {
			isa = PBXGroup;
			children = (
				13FF5FCF1D9861E3008C7DFB /* Extensions */,
				13FF5FCC1D9859DB008C7DFB /* Log */,
				1302F7421D95BA35005CE678 /* Bash */,
				131ACE741D9323C1007E71E9 /* Version */,
				13CA73DC1D84B63E00B1A323 /* Models */,
				13E3F5501D8342DD00AE7C20 /* Managers */,
				13E3F54D1D8341DF00AE7C20 /* Controllers */,
				13E3F54E1D83429900AE7C20 /* Supporting Files */,
				13E3F54F1D8342BC00AE7C20 /* Assets */,
				13C989801D8319600028BA2C /* AppDelegate.swift */,
			);
			path = BitriseStudio;
			sourceTree = "<group>";
		};
		13C989911D8319600028BA2C /* BitriseStudioTests */ = {
			isa = PBXGroup;
			children = (
				13C989921D8319600028BA2C /* BitriseStudioTests.swift */,
				13C989941D8319600028BA2C /* Info.plist */,
			);
			path = BitriseStudioTests;
			sourceTree = "<group>";
		};
		13C9899C1D8319600028BA2C /* BitriseStudioUITests */ = {
			isa = PBXGroup;
			children = (
				13C9899D1D8319600028BA2C /* BitriseStudioUITests.swift */,
				13C9899F1D8319600028BA2C /* Info.plist */,
			);
			path = BitriseStudioUITests;
			sourceTree = "<group>";
		};
		13CA73D71D84B5C500B1A323 /* Frameworks */ = {
			isa = PBXGroup;
			children = (
				131A3DA11D90609E002DAF99 /* Yaml.framework */,
				131A3D961D9053BB002DAF99 /* Realm.framework */,
				131A3D971D9053BB002DAF99 /* RealmSwift.framework */,
			);
			name = Frameworks;
			sourceTree = "<group>";
		};
		13CA73DC1D84B63E00B1A323 /* Models */ = {
			isa = PBXGroup;
			children = (
				13CA73DD1D84B67200B1A323 /* Project.swift */,
				131A3DAA1D906BFA002DAF99 /* BitriseTool.swift */,
				131A3DAC1D906D54002DAF99 /* Bitrise.swift */,
				131A3DAE1D906DAB002DAF99 /* Envman.swift */,
				131A3DB01D906DBF002DAF99 /* Stepman.swift */,
			);
			name = Models;
			sourceTree = "<group>";
		};
		13E3F54D1D8341DF00AE7C20 /* Controllers */ = {
			isa = PBXGroup;
			children = (
				13C989861D8319600028BA2C /* Main.storyboard */,
				13C989821D8319600028BA2C /* ProjectsViewController.swift */,
				13CA73D51D84A74800B1A323 /* AddProjectViewController.swift */,
				13B958CC1D89E87600D3310D /* SystemInfoViewController.swift */,
				13B958CE1D89EC7800D3310D /* RunViewController.swift */,
				13A95E811D966F040061B54F /* BashSessionViewController.swift */,
			);
			name = Controllers;
			sourceTree = "<group>";
		};
		13E3F54E1D83429900AE7C20 /* Supporting Files */ = {
			isa = PBXGroup;
			children = (
				13C989891D8319600028BA2C /* Info.plist */,
			);
			name = "Supporting Files";
			sourceTree = "<group>";
		};
		13E3F54F1D8342BC00AE7C20 /* Assets */ = {
			isa = PBXGroup;
			children = (
				13C989841D8319600028BA2C /* Assets.xcassets */,
			);
			name = Assets;
			sourceTree = "<group>";
		};
		13E3F5501D8342DD00AE7C20 /* Managers */ = {
			isa = PBXGroup;
			children = (
				13AAE1BC1D8426FF00AEE66D /* FileManager.swift */,
				13B0B8F21D872E93006EA29C /* RealmManager.swift */,
				131ACE721D93054B007E71E9 /* ToolManager.swift */,
			);
			name = Managers;
			sourceTree = "<group>";
		};
		13FF5FCC1D9859DB008C7DFB /* Log */ = {
			isa = PBXGroup;
			children = (
				13FF5FCD1D9859EE008C7DFB /* Log.swift */,
			);
			name = Log;
			sourceTree = "<group>";
		};
		13FF5FCF1D9861E3008C7DFB /* Extensions */ = {
			isa = PBXGroup;
			children = (
				13FF5FD01D98620A008C7DFB /* String+Extensions.swift */,
				13FF5FD21D9862DA008C7DFB /* Data+Extensions.swift */,
				13FF5FD41D9872EC008C7DFB /* Pipe+Extensions.swift */,
			);
			name = Extensions;
			sourceTree = "<group>";
		};
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
		13C9897C1D83195F0028BA2C /* BitriseStudio */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13C989A21D8319600028BA2C /* Build configuration list for PBXNativeTarget "BitriseStudio" */;
			buildPhases = (
				13C989791D83195F0028BA2C /* Sources */,
				13C9897A1D83195F0028BA2C /* Frameworks */,
				13C9897B1D83195F0028BA2C /* Resources */,
				131A3D9E1D90543F002DAF99 /* Embed Frameworks */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = BitriseStudio;
			productName = BitriseStudio;
			productReference = 13C9897D1D8319600028BA2C /* BitriseStudio.app */;
			productType = "com.apple.product-type.application";
		};
		13C9898D1D8319600028BA2C /* BitriseStudioTests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13C989A51D8319600028BA2C /* Build configuration list for PBXNativeTarget "BitriseStudioTests" */;
			buildPhases = (
				13C9898A1D8319600028BA2C /* Sources */,
				13C9898B1D8319600028BA2C /* Frameworks */,
				13C9898C1D8319600028BA2C /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
				13C989901D8319600028BA2C /* PBXTargetDependency */,
			);
			name = BitriseStudioTests;
			productName = BitriseStudioTests;
			productReference = 13C9898E1D8319600028BA2C /* BitriseStudioTests.xctest */;
			productType = "com.apple.product-type.bundle.unit-test";
		};
		13C989981D8319600028BA2C /* BitriseStudioUITests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13C989A81D8319600028BA2C /* Build configuration list for PBXNativeTarget "BitriseStudioUITests" */;
			buildPhases = (
				13C989951D8319600028BA2C /* Sources */,
				13C989961D8319600028BA2C /* Frameworks */,
				13C989971D8319600028BA2C /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
				13C9899B1D8319600028BA2C /* PBXTargetDependency */,
			);
			name = BitriseStudioUITests;
			productName = BitriseStudioUITests;
			productReference = 13C989991D8319600028BA2C /* BitriseStudioUITests.xctest */;
			productType = "com.apple.product-type.bundle.ui-testing";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		13C989751D83195F0028BA2C /* Project object */ = {
			isa = PBXProject;
			attributes = {
				LastSwiftUpdateCheck = 0800;
				LastUpgradeCheck = 0810;
				ORGANIZATIONNAME = "Krisztian Goedrei";
				TargetAttributes = {
					13C9897C1D83195F0028BA2C = {
						CreatedOnToolsVersion = 8.0;
						DevelopmentTeam = 9NS44DLTN7;
						ProvisioningStyle = Manual;
					};
					13C9898D1D8319600028BA2C = {
						CreatedOnToolsVersion = 8.0;
						DevelopmentTeam = L935L4GU3F;
						ProvisioningStyle = Automatic;
						TestTargetID = 13C9897C1D83195F0028BA2C;
					};
					13C989981D8319600028BA2C = {
						CreatedOnToolsVersion = 8.0;
						DevelopmentTeam = L935L4GU3F;
						ProvisioningStyle = Automatic;
						TestTargetID = 13C9897C1D83195F0028BA2C;
					};
				};
			};
			buildConfigurationList = 13C989781D83195F0028BA2C /* Build configuration list for PBXProject "BitriseStudio" */;
			compatibilityVersion = "Xcode 3.2";
			developmentRegion = English;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 13C989741D83195F0028BA2C;
			productRefGroup = 13C9897E1D8319600028BA2C /* Products */;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				13C9897C1D83195F0028BA2C /* BitriseStudio */,
				13C9898D1D8319600028BA2C /* BitriseStudioTests */,
				13C989981D8319600028BA2C /* BitriseStudioUITests */,
			);
		};
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
		13C9897B1D83195F0028BA2C /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C989851D8319600028BA2C /* Assets.xcassets in Resources */,
				13C989881D8319600028BA2C /* Main.storyboard in Resources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C9898C1D8319600028BA2C /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C989971D8319600028BA2C /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXResourcesBuildPhase section */

/* Begin PBXSourcesBuildPhase section */
		13C989791D83195F0028BA2C /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13FF5FD11D98620A008C7DFB /* String+Extensions.swift in Sources */,
				131A3DB11D906DBF002DAF99 /* Stepman.swift in Sources */,
				13B0B8F31D872E93006EA29C /* RealmManager.swift in Sources */,
				13FF5FCE1D9859EE008C7DFB /* Log.swift in Sources */,
				131A3DAB1D906BFA002DAF99 /* BitriseTool.swift in Sources */,
				13E3F5531D83477300AE7C20 /* ProjectsViewController.swift in Sources */,
				130E6BBE1D95BBB4009D3C78 /* Command.swift in Sources */,
				13FF5FD51D9872EC008C7DFB /* Pipe+Extensions.swift in Sources */,
				131ACE731D93054B007E71E9 /* ToolManager.swift in Sources */,
				13A95E821D966F040061B54F /* BashSessionViewController.swift in Sources */,
				13C989811D8319600028BA2C /* AppDelegate.swift in Sources */,
				1302F7441D95BA4A005CE678 /* Session.swift in Sources */,
				131ACE761D9323E3007E71E9 /* Version.swift in Sources */,
				13B958CD1D89E87600D3310D /* SystemInfoViewController.swift in Sources */,
				13CA73D61D84A74800B1A323 /* AddProjectViewController.swift in Sources */,
				131A3DAD1D906D54002DAF99 /* Bitrise.swift in Sources */,
				13FF5FD31D9862DA008C7DFB /* Data+Extensions.swift in Sources */,
				13AAE1BD1D8426FF00AEE66D /* FileManager.swift in Sources */,
				13B958D01D89EC7800D3310D /* RunViewController.swift in Sources */,
				131A3DAF1D906DAB002DAF99 /* Envman.swift in Sources */,
				13CA73DE1D84B67200B1A323 /* Project.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C9898A1D8319600028BA2C /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C989931D8319600028BA2C /* BitriseStudioTests.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C989951D8319600028BA2C /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C9899E1D8319600028BA2C /* BitriseStudioUITests.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXSourcesBuildPhase section */

/* Begin PBXTargetDependency section */
		13C989901D8319600028BA2C /* PBXTargetDependency */ = {
			isa = PBXTargetDependency;
			target = 13C9897C1D83195F0028BA2C /* BitriseStudio */;
			targetProxy = 13C9898F1D8319600028BA2C /* PBXContainerItemProxy */;
		};
		13C9899B1D8319600028BA2C /* PBXTargetDependency */ = {
			isa = PBXTargetDependency;
			target = 13C9897C1D83195F0028BA2C /* BitriseStudio */;
			targetProxy = 13C9899A1D8319600028BA2C /* PBXContainerItemProxy */;
		};
/* End PBXTargetDependency section */

/* Begin PBXVariantGroup section */
		13C989861D8319600028BA2C /* Main.storyboard */ = {
			isa = PBXVariantGroup;
			children = (
				13C989871D8319600028BA2C /* Base */,
			);
			name = Main.storyboard;
			sourceTree = "<group>";
		};
/* End PBXVariantGroup section */

/* Begin XCBuildConfiguration section */
		13C989A01D8319600028BA2C /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_ANALYZER_NONNULL = YES;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_DOCUMENTATION_COMMENTS = YES;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INFINITE_RECURSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN_SUSPICIOUS_MOVE = YES;
				CLANG_WARN_SUSPICIOUS_MOVES = YES;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				CODE_SIGN_IDENTITY = "-";
				COPY_PHASE_STRIP = NO;
				DEBUG_INFORMATION_FORMAT = dwarf;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				ENABLE_TESTABILITY = YES;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_DYNAMIC_NO_PIC = NO;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_OPTIMIZATION_LEVEL = 0;
				GCC_PREPROCESSOR_DEFINITIONS = (
					"DEBUG=1",
					"$(inherited)",
				);
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				MACOSX_DEPLOYMENT_TARGET = 10.11;
				MTL_ENABLE_DEBUG_INFO = YES;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = macosx;
				SWIFT_ACTIVE_COMPILATION_CONDITIONS = DEBUG;
				SWIFT_OPTIMIZATION_LEVEL = "-Onone";
			};
			name = Debug;
		};
		13C989A11D8319600028BA2C /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_ANALYZER_NONNULL = YES;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_DOCUMENTATION_COMMENTS = YES;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INFINITE_RECURSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN_SUSPICIOUS_MOVE = YES;
				CLANG_WARN_SUSPICIOUS_MOVES = YES;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				CODE_SIGN_IDENTITY = "-";
				COPY_PHASE_STRIP = NO;
				DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
				ENABLE_NS_ASSERTIONS = NO;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				MACOSX_DEPLOYMENT_TARGET = 10.11;
				MTL_ENABLE_DEBUG_INFO = NO;
				SDKROOT = macosx;
				SWIFT_OPTIMIZATION_LEVEL = "-Owholemodule";
			};
			name = Release;
		};
		13C989A31D8319600028BA2C /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_IDENTITY = "Mac Developer: Some Dude (KYXQXCWE3G)";
				COMBINE_HIDPI_IMAGES = YES;
				DEVELOPMENT_TEAM = 9NS44DLTN7;
				FRAMEWORK_SEARCH_PATHS = (
					"$(inherited)",
					"$(PROJECT_DIR)/Carthage/Build/Mac",
					"$(PROJECT_DIR)/Framworks",
				);
				INFOPLIST_FILE = BitriseStudio/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/../Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseStudio;
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE = "b17a1b90-9459-4620-9332-347d399f7cd9";
				PROVISIONING_PROFILE_SPECIFIER = "Mac Development Wildcard";
				SWIFT_VERSION = 3.0;
			};
			name = Debug;
		};
		13C989A41D8319600028BA2C /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				CODE_SIGN_IDENTITY = "3rd Party Mac Developer Application";
				COMBINE_HIDPI_IMAGES = YES;
				DEVELOPMENT_TEAM = 9NS44DLTN7;
				FRAMEWORK_SEARCH_PATHS = (
					"$(inherited)",
					"$(PROJECT_DIR)/Carthage/Build/Mac",
					"$(PROJECT_DIR)/Framworks",
				);
				INFOPLIST_FILE = BitriseStudio/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/../Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseStudio;
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE = "1bb807b8-a953-459e-85ca-c86d3fe13645";
				PROVISIONING_PROFILE_SPECIFIER = "Mac App-Store Wildcards";
				SWIFT_VERSION = 3.0;
			};
			name = Release;
		};
		13C989A61D8319600028BA2C /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				BUNDLE_LOADER = "$(TEST_HOST)";
				COMBINE_HIDPI_IMAGES = YES;
				DEVELOPMENT_TEAM = L935L4GU3F;
				INFOPLIST_FILE = BitriseStudioTests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/../Frameworks @loader_path/../Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = bitrise.BitriseStudioTests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_HOST = "$(BUILT_PRODUCTS_DIR)/BitriseStudio.app/Contents/MacOS/BitriseStudio";
			};
			name = Debug;
		};
		13C989A71D8319600028BA2C /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				BUNDLE_LOADER = "$(TEST_HOST)";
				COMBINE_HIDPI_IMAGES = YES;
				DEVELOPMENT_TEAM = L935L4GU3F;
				INFOPLIST_FILE = BitriseStudioTests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/../Frameworks @loader_path/../Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = bitrise.BitriseStudioTests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_HOST = "$(BUILT_PRODUCTS_DIR)/BitriseStudio.app/Contents/MacOS/BitriseStudio";
			};
			name = Release;
		};
		13C989A91D8319600028BA2C /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				COMBINE_HIDPI_IMAGES = YES;
				DEVELOPMENT_TEAM = L935L4GU3F;
				INFOPLIST_FILE = BitriseStudioUITests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/../Frameworks @loader_path/../Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = bitrise.BitriseStudioUITests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_TARGET_NAME = BitriseStudio;
			};
			name = Debug;
		};
		13C989AA1D8319600028BA2C /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				COMBINE_HIDPI_IMAGES = YES;
				DEVELOPMENT_TEAM = L935L4GU3F;
				INFOPLIST_FILE = BitriseStudioUITests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/../Frameworks @loader_path/../Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = bitrise.BitriseStudioUITests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_TARGET_NAME = BitriseStudio;
			};
			name = Release;
		};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		13C989781D83195F0028BA2C /* Build configuration list for PBXProject "BitriseStudio" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C989A01D8319600028BA2C /* Debug */,
				13C989A11D8319600028BA2C /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13C989A21D8319600028BA2C /* Build configuration list for PBXNativeTarget "BitriseStudio" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C989A31D8319600028BA2C /* Debug */,
				13C989A41D8319600028BA2C /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13C989A51D8319600028BA2C /* Build configuration list for PBXNativeTarget "BitriseStudioTests" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C989A61D8319600028BA2C /* Debug */,
				13C989A71D8319600028BA2C /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13C989A81D8319600028BA2C /* Build configuration list for PBXNativeTarget "BitriseStudioUITests" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C989A91D8319600028BA2C /* Debug */,
				13C989AA1D8319600028BA2C /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
/* End XCConfigurationList section */
	};
	rootObject = 13C989751D83195F0028BA2C /* Project object */;
}
`

const testIOSPbxprojContent = `
// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 46;
	objects = {

/* Begin PBXBuildFile section */
		13C4D5AB1DDDDED300D5DC29 /* AppDelegate.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C4D5AA1DDDDED300D5DC29 /* AppDelegate.swift */; };
		13C4D5AD1DDDDED300D5DC29 /* ViewController.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C4D5AC1DDDDED300D5DC29 /* ViewController.swift */; };
		13C4D5B01DDDDED300D5DC29 /* Main.storyboard in Resources */ = {isa = PBXBuildFile; fileRef = 13C4D5AE1DDDDED300D5DC29 /* Main.storyboard */; };
		13C4D5B21DDDDED300D5DC29 /* Assets.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 13C4D5B11DDDDED300D5DC29 /* Assets.xcassets */; };
		13C4D5B51DDDDED300D5DC29 /* LaunchScreen.storyboard in Resources */ = {isa = PBXBuildFile; fileRef = 13C4D5B31DDDDED300D5DC29 /* LaunchScreen.storyboard */; };
		13C4D5C01DDDDED400D5DC29 /* BitriseFastlaneSampleTests.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C4D5BF1DDDDED400D5DC29 /* BitriseFastlaneSampleTests.swift */; };
		13C4D5CB1DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.swift in Sources */ = {isa = PBXBuildFile; fileRef = 13C4D5CA1DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.swift */; };
/* End PBXBuildFile section */

/* Begin PBXContainerItemProxy section */
		13C4D5BC1DDDDED400D5DC29 /* PBXContainerItemProxy */ = {
			isa = PBXContainerItemProxy;
			containerPortal = 13C4D59F1DDDDED300D5DC29 /* Project object */;
			proxyType = 1;
			remoteGlobalIDString = 13C4D5A61DDDDED300D5DC29;
			remoteInfo = BitriseFastlaneSample;
		};
		13C4D5C71DDDDED400D5DC29 /* PBXContainerItemProxy */ = {
			isa = PBXContainerItemProxy;
			containerPortal = 13C4D59F1DDDDED300D5DC29 /* Project object */;
			proxyType = 1;
			remoteGlobalIDString = 13C4D5A61DDDDED300D5DC29;
			remoteInfo = BitriseFastlaneSample;
		};
/* End PBXContainerItemProxy section */

/* Begin PBXFileReference section */
		13C4D5A71DDDDED300D5DC29 /* BitriseFastlaneSample.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = BitriseFastlaneSample.app; sourceTree = BUILT_PRODUCTS_DIR; };
		13C4D5AA1DDDDED300D5DC29 /* AppDelegate.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = AppDelegate.swift; sourceTree = "<group>"; };
		13C4D5AC1DDDDED300D5DC29 /* ViewController.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = ViewController.swift; sourceTree = "<group>"; };
		13C4D5AF1DDDDED300D5DC29 /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/Main.storyboard; sourceTree = "<group>"; };
		13C4D5B11DDDDED300D5DC29 /* Assets.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Assets.xcassets; sourceTree = "<group>"; };
		13C4D5B41DDDDED300D5DC29 /* Base */ = {isa = PBXFileReference; lastKnownFileType = file.storyboard; name = Base; path = Base.lproj/LaunchScreen.storyboard; sourceTree = "<group>"; };
		13C4D5B61DDDDED300D5DC29 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13C4D5BB1DDDDED400D5DC29 /* BitriseFastlaneSampleTests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = BitriseFastlaneSampleTests.xctest; sourceTree = BUILT_PRODUCTS_DIR; };
		13C4D5BF1DDDDED400D5DC29 /* BitriseFastlaneSampleTests.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = BitriseFastlaneSampleTests.swift; sourceTree = "<group>"; };
		13C4D5C11DDDDED400D5DC29 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
		13C4D5C61DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.xctest */ = {isa = PBXFileReference; explicitFileType = wrapper.cfbundle; includeInIndex = 0; path = BitriseFastlaneSampleUITests.xctest; sourceTree = BUILT_PRODUCTS_DIR; };
		13C4D5CA1DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.swift */ = {isa = PBXFileReference; lastKnownFileType = sourcecode.swift; path = BitriseFastlaneSampleUITests.swift; sourceTree = "<group>"; };
		13C4D5CC1DDDDED400D5DC29 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXFrameworksBuildPhase section */
		13C4D5A41DDDDED300D5DC29 /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C4D5B81DDDDED400D5DC29 /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C4D5C31DDDDED400D5DC29 /* Frameworks */ = {
			isa = PBXFrameworksBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXFrameworksBuildPhase section */

/* Begin PBXGroup section */
		13C4D59E1DDDDED300D5DC29 = {
			isa = PBXGroup;
			children = (
				13C4D5A91DDDDED300D5DC29 /* BitriseFastlaneSample */,
				13C4D5BE1DDDDED400D5DC29 /* BitriseFastlaneSampleTests */,
				13C4D5C91DDDDED400D5DC29 /* BitriseFastlaneSampleUITests */,
				13C4D5A81DDDDED300D5DC29 /* Products */,
			);
			sourceTree = "<group>";
		};
		13C4D5A81DDDDED300D5DC29 /* Products */ = {
			isa = PBXGroup;
			children = (
				13C4D5A71DDDDED300D5DC29 /* BitriseFastlaneSample.app */,
				13C4D5BB1DDDDED400D5DC29 /* BitriseFastlaneSampleTests.xctest */,
				13C4D5C61DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.xctest */,
			);
			name = Products;
			sourceTree = "<group>";
		};
		13C4D5A91DDDDED300D5DC29 /* BitriseFastlaneSample */ = {
			isa = PBXGroup;
			children = (
				13C4D5AA1DDDDED300D5DC29 /* AppDelegate.swift */,
				13C4D5AC1DDDDED300D5DC29 /* ViewController.swift */,
				13C4D5AE1DDDDED300D5DC29 /* Main.storyboard */,
				13C4D5B11DDDDED300D5DC29 /* Assets.xcassets */,
				13C4D5B31DDDDED300D5DC29 /* LaunchScreen.storyboard */,
				13C4D5B61DDDDED300D5DC29 /* Info.plist */,
			);
			path = BitriseFastlaneSample;
			sourceTree = "<group>";
		};
		13C4D5BE1DDDDED400D5DC29 /* BitriseFastlaneSampleTests */ = {
			isa = PBXGroup;
			children = (
				13C4D5BF1DDDDED400D5DC29 /* BitriseFastlaneSampleTests.swift */,
				13C4D5C11DDDDED400D5DC29 /* Info.plist */,
			);
			path = BitriseFastlaneSampleTests;
			sourceTree = "<group>";
		};
		13C4D5C91DDDDED400D5DC29 /* BitriseFastlaneSampleUITests */ = {
			isa = PBXGroup;
			children = (
				13C4D5CA1DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.swift */,
				13C4D5CC1DDDDED400D5DC29 /* Info.plist */,
			);
			path = BitriseFastlaneSampleUITests;
			sourceTree = "<group>";
		};
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
		13C4D5A61DDDDED300D5DC29 /* BitriseFastlaneSample */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13C4D5CF1DDDDED400D5DC29 /* Build configuration list for PBXNativeTarget "BitriseFastlaneSample" */;
			buildPhases = (
				13C4D5A31DDDDED300D5DC29 /* Sources */,
				13C4D5A41DDDDED300D5DC29 /* Frameworks */,
				13C4D5A51DDDDED300D5DC29 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
			);
			name = BitriseFastlaneSample;
			productName = BitriseFastlaneSample;
			productReference = 13C4D5A71DDDDED300D5DC29 /* BitriseFastlaneSample.app */;
			productType = "com.apple.product-type.application";
		};
		13C4D5BA1DDDDED400D5DC29 /* BitriseFastlaneSampleTests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13C4D5D21DDDDED400D5DC29 /* Build configuration list for PBXNativeTarget "BitriseFastlaneSampleTests" */;
			buildPhases = (
				13C4D5B71DDDDED400D5DC29 /* Sources */,
				13C4D5B81DDDDED400D5DC29 /* Frameworks */,
				13C4D5B91DDDDED400D5DC29 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
				13C4D5BD1DDDDED400D5DC29 /* PBXTargetDependency */,
			);
			name = BitriseFastlaneSampleTests;
			productName = BitriseFastlaneSampleTests;
			productReference = 13C4D5BB1DDDDED400D5DC29 /* BitriseFastlaneSampleTests.xctest */;
			productType = "com.apple.product-type.bundle.unit-test";
		};
		13C4D5C51DDDDED400D5DC29 /* BitriseFastlaneSampleUITests */ = {
			isa = PBXNativeTarget;
			buildConfigurationList = 13C4D5D51DDDDED400D5DC29 /* Build configuration list for PBXNativeTarget "BitriseFastlaneSampleUITests" */;
			buildPhases = (
				13C4D5C21DDDDED400D5DC29 /* Sources */,
				13C4D5C31DDDDED400D5DC29 /* Frameworks */,
				13C4D5C41DDDDED400D5DC29 /* Resources */,
			);
			buildRules = (
			);
			dependencies = (
				13C4D5C81DDDDED400D5DC29 /* PBXTargetDependency */,
			);
			name = BitriseFastlaneSampleUITests;
			productName = BitriseFastlaneSampleUITests;
			productReference = 13C4D5C61DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.xctest */;
			productType = "com.apple.product-type.bundle.ui-testing";
		};
/* End PBXNativeTarget section */

/* Begin PBXProject section */
		13C4D59F1DDDDED300D5DC29 /* Project object */ = {
			isa = PBXProject;
			attributes = {
				LastSwiftUpdateCheck = 0810;
				LastUpgradeCheck = 0810;
				ORGANIZATIONNAME = "Krisztian Goedrei";
				TargetAttributes = {
					13C4D5A61DDDDED300D5DC29 = {
						CreatedOnToolsVersion = 8.1;
						DevelopmentTeam = 9NS44DLTN7;
						ProvisioningStyle = Manual;
					};
					13C4D5BA1DDDDED400D5DC29 = {
						CreatedOnToolsVersion = 8.1;
						DevelopmentTeam = 72SA8V3WYL;
						ProvisioningStyle = Automatic;
						TestTargetID = 13C4D5A61DDDDED300D5DC29;
					};
					13C4D5C51DDDDED400D5DC29 = {
						CreatedOnToolsVersion = 8.1;
						DevelopmentTeam = 72SA8V3WYL;
						ProvisioningStyle = Automatic;
						TestTargetID = 13C4D5A61DDDDED300D5DC29;
					};
				};
			};
			buildConfigurationList = 13C4D5A21DDDDED300D5DC29 /* Build configuration list for PBXProject "BitriseFastlaneSample" */;
			compatibilityVersion = "Xcode 3.2";
			developmentRegion = English;
			hasScannedForEncodings = 0;
			knownRegions = (
				en,
				Base,
			);
			mainGroup = 13C4D59E1DDDDED300D5DC29;
			productRefGroup = 13C4D5A81DDDDED300D5DC29 /* Products */;
			projectDirPath = "";
			projectRoot = "";
			targets = (
				13C4D5A61DDDDED300D5DC29 /* BitriseFastlaneSample */,
				13C4D5BA1DDDDED400D5DC29 /* BitriseFastlaneSampleTests */,
				13C4D5C51DDDDED400D5DC29 /* BitriseFastlaneSampleUITests */,
			);
		};
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
		13C4D5A51DDDDED300D5DC29 /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C4D5B51DDDDED300D5DC29 /* LaunchScreen.storyboard in Resources */,
				13C4D5B21DDDDED300D5DC29 /* Assets.xcassets in Resources */,
				13C4D5B01DDDDED300D5DC29 /* Main.storyboard in Resources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C4D5B91DDDDED400D5DC29 /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C4D5C41DDDDED400D5DC29 /* Resources */ = {
			isa = PBXResourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXResourcesBuildPhase section */

/* Begin PBXSourcesBuildPhase section */
		13C4D5A31DDDDED300D5DC29 /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C4D5AD1DDDDED300D5DC29 /* ViewController.swift in Sources */,
				13C4D5AB1DDDDED300D5DC29 /* AppDelegate.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C4D5B71DDDDED400D5DC29 /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C4D5C01DDDDED400D5DC29 /* BitriseFastlaneSampleTests.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
		13C4D5C21DDDDED400D5DC29 /* Sources */ = {
			isa = PBXSourcesBuildPhase;
			buildActionMask = 2147483647;
			files = (
				13C4D5CB1DDDDED400D5DC29 /* BitriseFastlaneSampleUITests.swift in Sources */,
			);
			runOnlyForDeploymentPostprocessing = 0;
		};
/* End PBXSourcesBuildPhase section */

/* Begin PBXTargetDependency section */
		13C4D5BD1DDDDED400D5DC29 /* PBXTargetDependency */ = {
			isa = PBXTargetDependency;
			target = 13C4D5A61DDDDED300D5DC29 /* BitriseFastlaneSample */;
			targetProxy = 13C4D5BC1DDDDED400D5DC29 /* PBXContainerItemProxy */;
		};
		13C4D5C81DDDDED400D5DC29 /* PBXTargetDependency */ = {
			isa = PBXTargetDependency;
			target = 13C4D5A61DDDDED300D5DC29 /* BitriseFastlaneSample */;
			targetProxy = 13C4D5C71DDDDED400D5DC29 /* PBXContainerItemProxy */;
		};
/* End PBXTargetDependency section */

/* Begin PBXVariantGroup section */
		13C4D5AE1DDDDED300D5DC29 /* Main.storyboard */ = {
			isa = PBXVariantGroup;
			children = (
				13C4D5AF1DDDDED300D5DC29 /* Base */,
			);
			name = Main.storyboard;
			sourceTree = "<group>";
		};
		13C4D5B31DDDDED300D5DC29 /* LaunchScreen.storyboard */ = {
			isa = PBXVariantGroup;
			children = (
				13C4D5B41DDDDED300D5DC29 /* Base */,
			);
			name = LaunchScreen.storyboard;
			sourceTree = "<group>";
		};
/* End PBXVariantGroup section */

/* Begin XCBuildConfiguration section */
		13C4D5CD1DDDDED400D5DC29 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_ANALYZER_NONNULL = YES;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_DOCUMENTATION_COMMENTS = YES;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INFINITE_RECURSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN_SUSPICIOUS_MOVES = YES;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
				COPY_PHASE_STRIP = NO;
				DEBUG_INFORMATION_FORMAT = dwarf;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				ENABLE_TESTABILITY = YES;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_DYNAMIC_NO_PIC = NO;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_OPTIMIZATION_LEVEL = 0;
				GCC_PREPROCESSOR_DEFINITIONS = (
					"DEBUG=1",
					"$(inherited)",
				);
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				IPHONEOS_DEPLOYMENT_TARGET = 10.1;
				MTL_ENABLE_DEBUG_INFO = YES;
				ONLY_ACTIVE_ARCH = YES;
				SDKROOT = iphoneos;
				SWIFT_ACTIVE_COMPILATION_CONDITIONS = DEBUG;
				SWIFT_OPTIMIZATION_LEVEL = "-Onone";
				TARGETED_DEVICE_FAMILY = "1,2";
			};
			name = Debug;
		};
		13C4D5CE1DDDDED400D5DC29 /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_SEARCH_USER_PATHS = NO;
				CLANG_ANALYZER_NONNULL = YES;
				CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
				CLANG_CXX_LIBRARY = "libc++";
				CLANG_ENABLE_MODULES = YES;
				CLANG_ENABLE_OBJC_ARC = YES;
				CLANG_WARN_BOOL_CONVERSION = YES;
				CLANG_WARN_CONSTANT_CONVERSION = YES;
				CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
				CLANG_WARN_DOCUMENTATION_COMMENTS = YES;
				CLANG_WARN_EMPTY_BODY = YES;
				CLANG_WARN_ENUM_CONVERSION = YES;
				CLANG_WARN_INFINITE_RECURSION = YES;
				CLANG_WARN_INT_CONVERSION = YES;
				CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
				CLANG_WARN_SUSPICIOUS_MOVES = YES;
				CLANG_WARN_UNREACHABLE_CODE = YES;
				CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
				COPY_PHASE_STRIP = NO;
				DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
				ENABLE_NS_ASSERTIONS = NO;
				ENABLE_STRICT_OBJC_MSGSEND = YES;
				GCC_C_LANGUAGE_STANDARD = gnu99;
				GCC_NO_COMMON_BLOCKS = YES;
				GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
				GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
				GCC_WARN_UNDECLARED_SELECTOR = YES;
				GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
				GCC_WARN_UNUSED_FUNCTION = YES;
				GCC_WARN_UNUSED_VARIABLE = YES;
				IPHONEOS_DEPLOYMENT_TARGET = 10.1;
				MTL_ENABLE_DEBUG_INFO = NO;
				SDKROOT = iphoneos;
				SWIFT_OPTIMIZATION_LEVEL = "-Owholemodule";
				TARGETED_DEVICE_FAMILY = "1,2";
				VALIDATE_PRODUCT = YES;
			};
			name = Release;
		};
		13C4D5D01DDDDED400D5DC29 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Distribution";
				DEVELOPMENT_TEAM = 9NS44DLTN7;
				INFOPLIST_FILE = BitriseFastlaneSample/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseFastlaneSample;
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE = "8e4701a8-01fb-4467-aad7-5a6c541795f0";
				PROVISIONING_PROFILE_SPECIFIER = "match AppStore com.bitrise.BitriseFastlaneSample";
				SWIFT_VERSION = 3.0;
			};
			name = Debug;
		};
		13C4D5D11DDDDED400D5DC29 /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
				"CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Distribution";
				DEVELOPMENT_TEAM = 9NS44DLTN7;
				INFOPLIST_FILE = BitriseFastlaneSample/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseFastlaneSample;
				PRODUCT_NAME = "$(TARGET_NAME)";
				PROVISIONING_PROFILE = "8e4701a8-01fb-4467-aad7-5a6c541795f0";
				PROVISIONING_PROFILE_SPECIFIER = "match AppStore com.bitrise.BitriseFastlaneSample";
				SWIFT_VERSION = 3.0;
			};
			name = Release;
		};
		13C4D5D31DDDDED400D5DC29 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				BUNDLE_LOADER = "$(TEST_HOST)";
				DEVELOPMENT_TEAM = 72SA8V3WYL;
				INFOPLIST_FILE = BitriseFastlaneSampleTests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks @loader_path/Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseFastlaneSampleTests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_HOST = "$(BUILT_PRODUCTS_DIR)/BitriseFastlaneSample.app/BitriseFastlaneSample";
			};
			name = Debug;
		};
		13C4D5D41DDDDED400D5DC29 /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				BUNDLE_LOADER = "$(TEST_HOST)";
				DEVELOPMENT_TEAM = 72SA8V3WYL;
				INFOPLIST_FILE = BitriseFastlaneSampleTests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks @loader_path/Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseFastlaneSampleTests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_HOST = "$(BUILT_PRODUCTS_DIR)/BitriseFastlaneSample.app/BitriseFastlaneSample";
			};
			name = Release;
		};
		13C4D5D61DDDDED400D5DC29 /* Debug */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				DEVELOPMENT_TEAM = 72SA8V3WYL;
				INFOPLIST_FILE = BitriseFastlaneSampleUITests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks @loader_path/Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseFastlaneSampleUITests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_TARGET_NAME = BitriseFastlaneSample;
			};
			name = Debug;
		};
		13C4D5D71DDDDED400D5DC29 /* Release */ = {
			isa = XCBuildConfiguration;
			buildSettings = {
				ALWAYS_EMBED_SWIFT_STANDARD_LIBRARIES = YES;
				DEVELOPMENT_TEAM = 72SA8V3WYL;
				INFOPLIST_FILE = BitriseFastlaneSampleUITests/Info.plist;
				LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks @loader_path/Frameworks";
				PRODUCT_BUNDLE_IDENTIFIER = com.bitrise.BitriseFastlaneSampleUITests;
				PRODUCT_NAME = "$(TARGET_NAME)";
				SWIFT_VERSION = 3.0;
				TEST_TARGET_NAME = BitriseFastlaneSample;
			};
			name = Release;
		};
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
		13C4D5A21DDDDED300D5DC29 /* Build configuration list for PBXProject "BitriseFastlaneSample" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C4D5CD1DDDDED400D5DC29 /* Debug */,
				13C4D5CE1DDDDED400D5DC29 /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13C4D5CF1DDDDED400D5DC29 /* Build configuration list for PBXNativeTarget "BitriseFastlaneSample" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C4D5D01DDDDED400D5DC29 /* Debug */,
				13C4D5D11DDDDED400D5DC29 /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13C4D5D21DDDDED400D5DC29 /* Build configuration list for PBXNativeTarget "BitriseFastlaneSampleTests" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C4D5D31DDDDED400D5DC29 /* Debug */,
				13C4D5D41DDDDED400D5DC29 /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
		13C4D5D51DDDDED400D5DC29 /* Build configuration list for PBXNativeTarget "BitriseFastlaneSampleUITests" */ = {
			isa = XCConfigurationList;
			buildConfigurations = (
				13C4D5D61DDDDED400D5DC29 /* Debug */,
				13C4D5D71DDDDED400D5DC29 /* Release */,
			);
			defaultConfigurationIsVisible = 0;
			defaultConfigurationName = Release;
		};
/* End XCConfigurationList section */
	};
	rootObject = 13C4D59F1DDDDED300D5DC29 /* Project object */;
}
`
