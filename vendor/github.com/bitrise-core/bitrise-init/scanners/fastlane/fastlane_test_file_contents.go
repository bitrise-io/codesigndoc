package fastlane

const iosTesFastfileContent = `fastlane_version "1.49.0"

default_platform :ios


platform :ios do
  before_all do
    # ENV["SLACK_URL"] = "https://hooks.slack.com/services/..."

  end

  desc "Runs all tests, archives app"
  lane :test do

    match(
      username: ENV['FASTLANE_USERNAME'],
      app_identifier: "io.bitrise.BitriseFastlaneSample",
      readonly: true,
      type: "appstore"
    )

    scan(
      scheme: "BitriseFastlaneSample",
      destination: "name=iPhone 5s,OS=latest",
      output_directory: ENV['BITRISE_DEPLOY_DIR']
    )

    gym(
      scheme: "BitriseFastlaneSample",
      configuration: "Release",
      output_directory: ENV['BITRISE_DEPLOY_DIR'],
      output_name: "BitriseFastlaneSample.ipa",
      use_legacy_build_api: "true"
    )

    crashlytics(
      crashlytics_path: "./Crashlytics.framework",
      api_token: ENV['CRASHLYTICS_API_TOKEN'],
      build_secret: ENV['CRASHLYTICS_BUILD_SECRET'],
      ipa_path: "#{ENV['BITRISE_DEPLOY_DIR']}/BitriseFastlaneSample.ipa"
    )

    snapshot
  end

  after_all do |lane|
    
  end

  error do |lane, exception|
    
  end
end`

const complexIosTestFastFileContent = `default_platform :ios

platform :ios do
  before_all do
     # Set project for commit_version_bump, which seems to get confused by projects in other folders
     ENV['FL_BUILD_NUMBER_PROJECT'] = "Wikipedia.xcodeproj"
     ensure_git_status_clean unless ENV['FL_NO_ENSURE_CLEAN']
  end

  desc "Runs linting (and eventually static analysis)"
  lane :analyze do
    xcodebuild(
      workspace: "Wikipedia.xcworkspace",
      scheme: "Wikipedia",
      configuration: "Debug",
      sdk: 'iphonesimulator',
      destination: 'platform=iOS Simulator,OS=9.3,name=iPhone 6',
      analyze: true
    )
  end

  desc "Runs tests, version, tag, and push to the beta branch"
  lane :testAndPushBeta do
    verifyTestPlatforms
    bumpAndTagBeta
    badge(dark: true, shield: get_badge_version_string, shield_no_resize: true)
    beta
  end

  desc "Runs tests, version, tag, and push to the beta branch"
  lane :submitAndPushToMaster do
    bumpAndTagRelease
    store
  end

  desc "Runs tests on the primary platforms and configurations"
  lane :verifyTestPlatforms do
    verify(
     sim_os: 8.4,
     scheme: "WikipediaRTL"
    )
    verify(
     sim_os: 8.4
    )
    verify(
     scheme: "WikipediaRTL"
    )
    verify(
     junit: true
    )
  end

  desc "Runs unit tests, optionally generating a JUnit report."
  lane :verify do |options|
    scheme = options[:scheme] || 'Wikipedia'
    sim_os = options[:sim_os] || '9.3'
    destination = "platform=iOS Simulator,name=iPhone 6,OS=#{sim_os}"
    opts = {
      :scheme =>  scheme,
      :workspace => 'Wikipedia.xcworkspace',
      :configuration => 'Debug',
      :destination => destination,
      :buildlog_path => './build',
      :output_directory => './build/reports',
      :output_style => 'basic',
      :code_coverage => true,
      :xcargs => "TRAVIS=#{ENV["TRAVIS"]}"
    }
    opts[:output_types] = options[:junit] ? "junit" : ""
    scan(opts)
  end

  desc "Increment the app version patch"
  lane :bumpPatch do
    increment_version_number(
      bump_type: "patch"
    )
  end

  desc "Increment the app version minor"
  lane :bumpMinor do
    increment_version_number(
      bump_type: "minor"
    )
  end

  desc "Increment the app version major"
  lane :bumpMajor do
    increment_version_number(
      bump_type: "major"
    )
  end

  desc "Increment the app's build number without committing the changes. Returns a string of the new, bumped version."
  lane :bump do |options|
    opt_build_num = options[:build_number] || ENV['BUILD_NUMBER']
    if opt_build_num then
      increment_build_number(build_number: opt_build_num.to_i)
    else
      increment_build_number(build_number: get_build_number)
    end
    get_version_number
  end

  desc "Increment the app's beta build number, add a tag, and push to the beta branch."
  lane :bumpAndTagBeta do |options|
    sh "git fetch"
    sh "git checkout develop"
    sh "git pull"
    sh "git checkout beta"
    sh "git merge develop"

    increment_build_number

    new_version = get_version_number
    commit_version_bump
    push_to_git_remote(
       local_branch: 'beta',  # optional, aliased by 'branch', default: 'master'
       remote_branch: 'beta', # optional, default is set to local_branch
    )

    tag_name = "betas/#{new_version}"
    add_git_tag(tag: tag_name)
    sh "git push origin --tags"

  end

  desc "Increment the app's build number, add a tag, and push to the master branch."
  lane :bumpAndTagRelease do |options|
    sh "git fetch"
    sh "git checkout release"
    sh "git pull"
    sh "git checkout master"
    sh "git merge release"

    increment_build_number(build_number: get_release_build_number)

    new_version = get_version_number
    commit_version_bump
    push_to_git_remote(
       local_branch: 'master',  # optional, aliased by 'branch', default: 'master'
       remote_branch: 'master', # optional, default is set to local_branch
    )

    tag_name = "releases/#{new_version}"
    add_git_tag(tag: tag_name)
    sh "git push origin --tags"

  end


  desc "Returns a default changelog."
  lane :default_changelog do
    changelog = changelog_from_git_commits(
        between: [ENV['GIT_PREVIOUS_SUCCESSFUL_COMMIT'] || "HEAD^^^^^", "HEAD"],
        pretty: "- %s"
    )
    # HAX: strip emoji from changelog
    changelog = changelog.sub(/[\u{1F300}-\u{1F6FF}]/, '')
    Actions.lane_context[SharedValues::FL_CHANGELOG] = changelog
    puts changelog
    changelog
  end

  desc "Submit a new **Wikipedia Beta** build to Apple TestFlight for internal testing."
  lane :beta do

    sigh(
      adhoc: false,
      force: true
    )

    gym(
      configuration: "Beta",
      scheme: "Wikipedia Beta"
    )

    # changelog = default_changelog

    hockey_beta_id = ENV["HOCKEY_BETA"]
    raise "Must specifiy HockeyApp public identifier." unless hockey_beta_id.length > 0

    hockey(
      public_identifier: hockey_beta_id,
      # notes: changelog,
      notify: '0', # Means do not notify
      status: '1', # Means do not make available for download
    )

    pilot(
      skip_submission: false,
      distribute_external: false
    )

  end

  desc "Submit a new App Store release candidate Apple TestFlight for internal testing."
  lane :store do
    sigh(
      adhoc: false,
      force: true
    )

    gym(
      configuration: "Release",
      scheme: "Wikipedia"
    )

    hockey_prod_id = ENV["HOCKEY_PRODUCTION"]
    raise "Must specifiy HockeyApp public identifier." unless hockey_prod_id.length > 0

    hockey(
      public_identifier: hockey_prod_id,
      notify: '0', # Means do not notify
      status: '1', # Means do not make available for download
    )

    pilot(
      skip_submission: false,
      distribute_external: false
    )

  end


  desc "Upload a developer build to Hockey."
  lane :dev do
    sigh(
      adhoc: true,
      # Fastlane has issues forcing AdHoc profiles
      force: false
    )

    # force iTunes file sharing to be enabled (normally disabled for release builds)
    ENV['WMF_FORCE_ITUNES_FILE_SHARING'] = '1'
    # force debug menu to be shown
    ENV['WMF_FORCE_DEBUG_MENU'] = '1'

    gym(
      configuration: "AdHoc",
      scheme: "Wikipedia AdHoc",
      # both of these flags are required for ad hoc
      export_method: 'ad-hoc',
      use_legacy_build_api: true
    )

    hockey(
      notes: default_changelog,
      notify: '2', # Notify all testers
      status: '2', # Make available for download
      release_type: '2' # 'alpha' release type
    )
  end
end`
