package xcodeproj

const codeSignInfoScriptContent = `require 'xcodeproj'
require 'json'

def workspace_contained_projects(workspace_pth)
  workspace = Xcodeproj::Workspace.new_from_xcworkspace(workspace_pth)
  workspace_dir = File.dirname(workspace_pth)
  project_paths = []
  workspace.file_references.each do |ref|
    pth = ref.path
    next unless File.extname(pth) == '.xcodeproj'
    next if pth.end_with?('Pods/Pods.xcodeproj')

    project_path = File.expand_path(pth, workspace_dir)
    project_paths << project_path
  end

  project_paths
end

def shared_scheme_path(project_or_workspace_pth, scheme_name)
  File.join(project_or_workspace_pth, 'xcshareddata', 'xcschemes', scheme_name + '.xcscheme')
end

def user_scheme_path(project_or_workspace_pth, scheme_name, user_name)
  File.join(project_or_workspace_pth, 'xcuserdata', user_name + '.xcuserdatad', 'xcschemes', scheme_name + '.xcscheme')
end

def read_scheme(project_or_workspace_pth, scheme_name, user_name)
  project_paths = [project_or_workspace_pth]
  if File.extname(project_or_workspace_pth) == '.xcworkspace'
    project_paths += workspace_contained_projects(project_or_workspace_pth)
  end

  project_paths.each do |project_path|
    scheme_pth = shared_scheme_path(project_path, scheme_name)
    scheme_pth = user_scheme_path(project_path, scheme_name, user_name) unless File.exist?(scheme_pth)
    next unless File.exist?(scheme_pth)

    scheme = Xcodeproj::XCScheme.new(scheme_pth)
    container_dir = File.dirname(project_path)

    return {
      scheme: scheme,
      container_dir: container_dir
    }
  end

  nil
end

def project_buildable_target_mapping(project_dir, scheme)
  build_action = scheme.build_action
  return nil unless build_action

  entries = build_action.entries || []
  return nil if entries.empty?

  entries = entries.select(&:build_for_archiving?) || []
  return nil if entries.empty?

  mapping = {}

  entries.each do |entry|
    buildable_references = entry.buildable_references || []
    next if buildable_references.empty?

    buildable_references = buildable_references.reject do |r|
      r.target_name.to_s.empty? || r.target_referenced_container.to_s.empty?
    end
    next if buildable_references.empty?

    buildable_reference = entry.buildable_references.first

    target_name = buildable_reference.target_name.to_s
    container = buildable_reference.target_referenced_container.to_s.sub(/^container:/, '')
    next if target_name.empty? || container.empty?

    project_pth = File.expand_path(container, project_dir)
    next unless File.exist?(project_pth)

    project = Xcodeproj::Project.open(project_pth)
    next unless project

    target = project.targets.find { |t| t.name == target_name }
    next unless target
    next unless runnable_target?(target)

    targets = mapping[project_pth] || []
    targets << target
    mapping[project_pth] = targets
  end

  mapping
end

def runnable_target?(target)
  return false unless target.is_a?(Xcodeproj::Project::Object::PBXNativeTarget)

  product_reference = target.product_reference
  return false unless product_reference

  product_reference.path.end_with?('.app', '.appex')
end

def collect_dependent_targets(target, dependent_targets)
  dependent_targets << target

  dependencies = target.dependencies || []
  return dependent_targets if dependencies.empty?

  dependencies.each do |dependency|
    dependent_target = dependency.target
    next unless dependent_target
    next unless runnable_target?(dependent_target)

    collect_dependent_targets(dependent_target, dependent_targets)
  end

  dependent_targets
end

def find_archive_action_build_configuration_name(scheme)
  archive_action = scheme.archive_action
  return nil unless archive_action

  archive_action.build_configuration
end

def read_scheme_target_mapping(project_or_workspace_pth, scheme_name, user_name)
  scheme_container_dir = read_scheme(project_or_workspace_pth, scheme_name, user_name)
  raise "project (#{project_or_workspace_pth}) does not contain scheme: #{scheme_name}" unless scheme_container_dir
  scheme = scheme_container_dir[:scheme]
  container_dir = scheme_container_dir[:container_dir]

  configuration = find_archive_action_build_configuration_name(scheme)

  target_mapping = project_buildable_target_mapping(container_dir, scheme) || []
  raise 'scheme does not contain buildable target' if target_mapping.empty?

  project_targets = {}
  target_mapping.each do |project_pth, targets|
    targets.each do |target|
      dependent_targets = []
      dependent_targets = collect_dependent_targets(target, dependent_targets)

      project_targets[project_pth] = dependent_targets.collect(&:name)
    end
  end
  raise 'failed to collect buildable targets' if project_targets.empty?

  {
    configuration: configuration,
    project_targets: project_targets
  }
end

begin
  project_path = ENV['project']
  scheme_name = ENV['scheme']
  user_name = ENV['user']

  raise 'missing project_path' if project_path.to_s.empty?
  raise 'missing scheme_name' if scheme_name.to_s.empty?
  raise 'missing user_name' if user_name.to_s.empty?

  mapping = read_scheme_target_mapping(project_path, scheme_name, user_name)
  result = {
    data: mapping
  }
  result_json = JSON.pretty_generate(result).to_s
  puts result_json
rescue => e
  error_message = e.to_s + "\n" + e.backtrace.join("\n")
  result = {
    error: error_message
  }
  result_json = result.to_json.to_s
  puts result_json
  exit(1)
end
`

const gemfileContent = `source "https://rubygems.org"
gem "xcodeproj"
gem "json"
`
