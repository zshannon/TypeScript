require 'yaml'
require 'tmpdir'
require 'fileutils'
require 'open3'
require 'minitest/autorun'

class ReleaseWorkflowTest < Minitest::Test
  ROOT = File.expand_path('..', __dir__)
  WORKFLOW = YAML.load_file(File.join(ROOT, '.github/workflows/release-public-module.yml'))

  def setup
    @dir = Dir.mktmpdir('typescript-release-test-')
    @upstream = File.join(@dir, 'upstream')
    @fork = File.join(@dir, 'fork')
    FileUtils.mkdir_p(@upstream)
    git(@upstream, 'init', '-b', 'main')
    git(@upstream, 'config', 'user.name', 'Test')
    git(@upstream, 'config', 'user.email', 'test@example.com')
  end

  def teardown
    FileUtils.remove_entry(@dir)
  end

  def git(dir, *args)
    out, status = Open3.capture2e('git', '-C', dir, *args)
    raise out unless status.success?
    out.strip
  end

  def fixture(version)
    FileUtils.mkdir_p(File.join(@upstream, 'tsc/internal/core'))
    File.write(File.join(@upstream, 'tsc/internal/core/version.go'), "package core\nvar version = \"#{version}\"\n")
    File.write(File.join(@upstream, 'tsc/go.mod'), "module github.com/microsoft/TypeScript/tsc\n\ngo 1.26\n")
    git(@upstream, 'add', '.')
    git(@upstream, 'commit', '-m', version)
    git(@upstream, 'tag', "v#{version}") unless version.include?('-')
    git(@dir, 'clone', @upstream, @fork)
    FileUtils.mkdir_p(File.join(@fork, 'scripts'))
    FileUtils.cp(File.join(ROOT, 'scripts/export-go.go'), File.join(@fork, 'scripts/export-go.go'))
    FileUtils.mkdir_p(File.join(@fork, 'public'))
    File.write(File.join(@fork, 'public/go.mod'), "module github.com/zshannon/TypeScript/public/v7\n")
  end

  def step(id, env = {})
    item = WORKFLOW.fetch('jobs').values.flat_map { |job| job.fetch('steps') }.find { |s| s['id'] == id || s['name'] == id }
    script = item.fetch('run').gsub('https://github.com/microsoft/TypeScript.git', @upstream)
    output = File.join(@dir, 'output')
    File.write(output, '')
    text, status = Open3.capture2e({'GOWORK' => 'off', 'GITHUB_OUTPUT' => output,
      'GITHUB_STEP_SUMMARY' => File.join(@dir, 'summary')}.merge(env), 'bash', '-c', script, chdir: @fork)
    [status.success?, text, File.read(output)]
  end

  def test_stable_version_and_upstream_tag
    fixture('7.0.2')
    ok, text, output = step('version')
    assert ok, text
    assert_includes output, 'public_tag=public/v7.0.2'
    ok, text = step('Verify the version against upstream source', 'SOURCE_VERSION' => 'v7.0.2')
    assert ok, text
  end

  def test_development_version_and_upstream_ancestor
    fixture('7.1.0-dev')
    ok, text, output = step('version')
    assert ok, text
    assert_includes output, 'public_tag=public/v7.1.0-dev'
    ok, text = step('Verify the version against upstream source', 'SOURCE_VERSION' => 'v7.1.0-dev')
    assert ok, text
  end

  def test_wrong_upstream_version_is_rejected
    fixture('7.1.0-dev')
    ok, = step('Verify the version against upstream source', 'SOURCE_VERSION' => 'v7.2.0-dev')
    refute ok
  end

  def test_wrong_module_major_is_rejected
    fixture('7.0.2')
    File.write(File.join(@fork, 'public/go.mod'), "module github.com/zshannon/TypeScript/public/v8\n")
    ok, = step('version')
    refute ok
  end

  def test_existing_version_is_skipped_and_new_version_is_released
    fixture('7.0.2')
    ok, text, output = step('existing', 'PUBLIC_TAG' => 'public/v7.0.2')
    assert ok, text
    assert_includes output, 'release_needed=true'
    git(@upstream, 'tag', 'public/v7.0.2')
    ok, text, output = step('existing', 'PUBLIC_TAG' => 'public/v7.0.2')
    assert ok, text
    assert_includes output, 'release_needed=false'
  end

  def test_publish_tags_tested_commit_and_never_retags
    fixture('7.0.2')
    original_sha = git(@fork, 'rev-parse', 'HEAD')
    env = {'PUBLIC_TAG' => 'public/v7.0.2', 'SOURCE_VERSION' => 'v7.0.2', 'TESTED_SHA' => original_sha}
    ok, text = step('Publish matching public module tag', env)
    assert ok, text
    assert_equal original_sha, git(@upstream, 'rev-parse', 'public/v7.0.2^{commit}')
    git(@fork, 'commit', '--allow-empty', '-m', 'Another merge with the same TypeScript version')
    env['TESTED_SHA'] = git(@fork, 'rev-parse', 'HEAD')
    ok, text = step('Publish matching public module tag', env)
    assert ok, text
    assert_equal original_sha, git(@upstream, 'rev-parse', 'public/v7.0.2^{commit}')
  end

  def test_remote_failure_does_not_publish
    fixture('7.0.2')
    git(@fork, 'remote', 'set-url', 'origin', File.join(@dir, 'missing'))
    ok, = step('existing', 'PUBLIC_TAG' => 'public/v7.0.2')
    refute ok
  end
end
