package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const (
	fixtureSourceModule = "github.com/microsoft/TypeScript/tsc"
	legacySourceModule  = "github.com/microsoft/typescript-go"
	fixturePublicModule = "github.com/zshannon/TypeScript/public/v7"
)

func TestExportModuleCopiesWholeTreeAndRewritesCompilerImports(t *testing.T) {
	repo := newFixtureRepo(t)

	if err := exportModule(repo, false); err != nil {
		t.Fatalf("exportModule() error = %v", err)
	}

	assertFileContent(t, filepath.Join(repo, "public", "go.mod"), "module "+fixturePublicModule+"\n\ngo 1.26\n")
	assertFileContent(t, filepath.Join(repo, "public", "go.sum"), "example.com/dependency v1.0.0 h1:fixture\n")
	assertFileContent(t, filepath.Join(repo, "public", "LICENSE"), "fixture license\n")
	assertFileContent(t, filepath.Join(repo, "public", "NOTICE"), "fixture notice\n")

	alpha := readFile(t, filepath.Join(repo, "public", "alpha", "alpha.go"))
	for _, want := range []string{
		`"` + fixturePublicModule + `/beta"`,
		`"` + fixturePublicModule + `/vfs/internals"`,
		"func Exported() {}",
		"func private() {}",
	} {
		if !bytes.Contains(alpha, []byte(want)) {
			t.Errorf("generated alpha.go does not contain %q:\n%s", want, alpha)
		}
	}
	if bytes.Contains(alpha, []byte(fixtureSourceModule+"/internal")) {
		t.Errorf("generated alpha.go retains source module import:\n%s", alpha)
	}
	generator := readFile(t, filepath.Join(repo, "public", "alpha", "generate.go"))
	if !bytes.Contains(generator, []byte(fixturePublicModule+"/vfs/internals")) {
		t.Errorf("generated source generator does not use the public nested package path:\n%s", generator)
	}
	if bytes.Contains(generator, []byte(fixtureSourceModule+"/internal")) {
		t.Errorf("generated source generator retains source module path:\n%s", generator)
	}

	assertFileContent(t, filepath.Join(repo, "public", "alpha", "tagged.go"), "//go:build fixturetag\n\npackage alpha\n")
	assertFileContent(t, filepath.Join(repo, "public", "vfs", "internals", "helper.go"), "package internal\n")
	assertFileContent(t, filepath.Join(repo, "public", "vfs", "internalized", "helper.go"), "package internalized\n")
	if _, err := os.Stat(filepath.Join(repo, "public", "vfs", "internal")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("source nested internal directory was not renamed: stat error = %v", err)
	}

	wantAsset := []byte{0, 1, 2, 0xff, '\n'}
	if got := readFile(t, filepath.Join(repo, "public", "assets", "data", "message.bin")); !bytes.Equal(got, wantAsset) {
		t.Errorf("binary asset = %v, want %v", got, wantAsset)
	}
	readme := readFile(t, filepath.Join(repo, "public", "README.md"))
	for _, want := range []string{
		"generated from every package under `tsc/internal`",
		fixturePublicModule + "/parser",
		fixturePublicModule + "@v7.0.2",
		"corresponds to TypeScript 7.0.2",
		"current generated source is `7.1.0-dev`, a development version",
		"go run ./scripts/export-go.go --check",
	} {
		if !bytes.Contains(readme, []byte(want)) {
			t.Errorf("generated README.md does not contain %q", want)
		}
	}
	if bytes.Contains(readme, []byte("pseudo-version")) {
		t.Error("generated README.md describes commit-only pseudo-version pinning")
	}
}

func TestExportModuleDerivesPublicMajorFromSourceVersion(t *testing.T) {
	repo := newFixtureRepo(t)
	if err := os.Remove(filepath.Join(repo, "NOTICE.txt")); err != nil {
		t.Fatalf("Remove(root NOTICE.txt): %v", err)
	}
	writeFixtureFile(t, filepath.Join(repo, "tsc", "NOTICE.txt"), []byte("compiler notice\n"))
	writeFixtureFile(t, filepath.Join(repo, "tsc", "internal", "core", "version.go"), []byte("package core\n\nvar version = \"8.2.3\"\n"))
	writeFixtureFile(t, filepath.Join(repo, "tsc", "go.mod"), []byte("module "+legacySourceModule+"\n\ngo 1.26\n"))
	writeFixtureFile(t, filepath.Join(repo, "tsc", "internal", "alpha", "alpha.go"), []byte("package alpha\n\nimport (\n\t_ \""+legacySourceModule+"/internal/beta\"\n\t_ \""+legacySourceModule+"/internal/vfs/internal\"\n)\n"))
	writeFixtureFile(t, filepath.Join(repo, "tsc", "internal", "alpha", "generate.go"), []byte("//go:build ignore\n\npackage main\n\nconst generatedImport = \"_ \\\""+legacySourceModule+"/internal/vfs/internal\\\"\"\n"))

	if err := exportModule(repo, false); err != nil {
		t.Fatalf("exportModule() error = %v", err)
	}

	const wantModule = "github.com/zshannon/TypeScript/public/v8"
	assertFileContent(t, filepath.Join(repo, "public", "go.mod"), "module "+wantModule+"\n\ngo 1.26\n")
	assertFileContent(t, filepath.Join(repo, "public", "NOTICE"), "compiler notice\n")
	alpha := readFile(t, filepath.Join(repo, "public", "alpha", "alpha.go"))
	if !bytes.Contains(alpha, []byte(wantModule+"/beta")) {
		t.Errorf("generated alpha.go does not use derived v8 module path:\n%s", alpha)
	}
	if bytes.Contains(alpha, []byte(legacySourceModule+"/internal")) {
		t.Errorf("generated alpha.go retains legacy source module path:\n%s", alpha)
	}
	readme := readFile(t, filepath.Join(repo, "public", "README.md"))
	if !bytes.Contains(readme, []byte(wantModule+"/parser")) {
		t.Errorf("generated README.md does not use derived v8 module path:\n%s", readme)
	}
	if !bytes.Contains(readme, []byte("current generated source is `8.2.3`, and its corresponding Go module version is `v8.2.3`")) {
		t.Errorf("generated README.md does not align the stable Go and TypeScript versions:\n%s", readme)
	}

	var stdout, stderr strings.Builder
	if exitCode := run([]string{"--version"}, repo, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("run(--version) exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	if got, want := stdout.String(), "v8.2.3\n"; got != want {
		t.Errorf("run(--version) output = %q, want %q", got, want)
	}
}

func TestExportModuleDeletesStaleFilesAndIsIdempotent(t *testing.T) {
	repo := newFixtureRepo(t)
	if err := exportModule(repo, false); err != nil {
		t.Fatalf("first exportModule() error = %v", err)
	}
	want := snapshotTree(t, filepath.Join(repo, "public"))

	writeFixtureFile(t, filepath.Join(repo, "public", "stale", "removed.txt"), []byte("stale\n"))
	if err := exportModule(repo, false); err != nil {
		t.Fatalf("second exportModule() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "public", "stale")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("stale directory remains after regeneration: stat error = %v", err)
	}
	if got := snapshotTree(t, filepath.Join(repo, "public")); !reflect.DeepEqual(got, want) {
		t.Errorf("second generation differs from first\ngot:  %#v\nwant: %#v", got, want)
	}
}

func TestCheckModeDetectsDriftWithoutWriting(t *testing.T) {
	repo := newFixtureRepo(t)
	if err := exportModule(repo, false); err != nil {
		t.Fatalf("exportModule() error = %v", err)
	}
	writeFixtureFile(t, filepath.Join(repo, "public", "alpha", "alpha.go"), []byte("modified\n"))
	writeFixtureFile(t, filepath.Join(repo, "public", "stale.txt"), []byte("stale\n"))
	wantUnchanged := snapshotTree(t, filepath.Join(repo, "public"))

	err := exportModule(repo, true)
	if !errors.Is(err, errOutputOutOfDate) {
		t.Fatalf("check exportModule() error = %v, want errOutputOutOfDate", err)
	}
	if got := snapshotTree(t, filepath.Join(repo, "public")); !reflect.DeepEqual(got, wantUnchanged) {
		t.Errorf("check mode changed output\ngot:  %#v\nwant: %#v", got, wantUnchanged)
	}

	if err := exportModule(repo, false); err != nil {
		t.Fatalf("repair exportModule() error = %v", err)
	}
	if err := exportModule(repo, true); err != nil {
		t.Fatalf("check after repair error = %v", err)
	}
}

func TestCheckModeDoesNotCreateMissingOutput(t *testing.T) {
	repo := newFixtureRepo(t)

	err := exportModule(repo, true)
	if !errors.Is(err, errOutputOutOfDate) {
		t.Fatalf("check exportModule() error = %v, want errOutputOutOfDate", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "public")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("check mode created output directory: stat error = %v", err)
	}
}

func newFixtureRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	files := map[string][]byte{
		"tsc/go.mod":                              []byte("module " + fixtureSourceModule + "\n\ngo 1.26\n"),
		"tsc/go.sum":                              []byte("example.com/dependency v1.0.0 h1:fixture\n"),
		"tsc/LICENSE":                             []byte("fixture license\n"),
		"NOTICE.txt":                              []byte("fixture notice\n"),
		"tsc/internal/alpha/alpha.go":             []byte("package alpha\n\nimport (\n\t_ \"" + fixtureSourceModule + "/internal/beta\"\n\t_ \"" + fixtureSourceModule + "/internal/vfs/internal\"\n)\n\nfunc Exported() {}\nfunc private() {}\n"),
		"tsc/internal/alpha/generate.go":          []byte("//go:build ignore\n\npackage main\n\nconst generatedImport = \"_ \\\"" + fixtureSourceModule + "/internal/vfs/internal\\\"\"\n"),
		"tsc/internal/alpha/tagged.go":            []byte("//go:build fixturetag\n\npackage alpha\n"),
		"tsc/internal/assets/embed.go":            []byte("package assets\n\nimport _ \"embed\"\n\n//go:embed data/message.bin\nvar message []byte\n"),
		"tsc/internal/assets/data/message.bin":    {0, 1, 2, 0xff, '\n'},
		"tsc/internal/beta/beta.go":               []byte("package beta\n"),
		"tsc/internal/core/version.go":            []byte("package core\n\nvar version = \"7.1.0-dev\"\n"),
		"tsc/internal/vfs/internal/helper.go":     []byte("package internal\n"),
		"tsc/internal/vfs/internalized/helper.go": []byte("package internalized\n"),
	}
	for name, content := range files {
		writeFixtureFile(t, filepath.Join(repo, filepath.FromSlash(name)), content)
	}
	return repo
}

func writeFixtureFile(t *testing.T, name string, content []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", filepath.Dir(name), err)
	}
	if err := os.WriteFile(name, content, 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", name, err)
	}
}

func assertFileContent(t *testing.T, name string, want string) {
	t.Helper()
	if got := string(readFile(t, name)); got != want {
		t.Errorf("%s content:\n%s\nwant:\n%s", name, got, want)
	}
}

func readFile(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", name, err)
	}
	return content
}

type treeEntry struct {
	Mode fs.FileMode
	Data string
}

func snapshotTree(t *testing.T, root string) map[string]treeEntry {
	t.Helper()
	result := map[string]treeEntry{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := treeEntry{Mode: info.Mode()}
		if !entry.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			item.Data = string(content)
		}
		result[filepath.ToSlash(rel)] = item
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(%q): %v", root, err)
	}
	return result
}

func TestRunAcceptsOnlyOptionalCheckFlag(t *testing.T) {
	repo := newFixtureRepo(t)
	var stdout, stderr strings.Builder
	if exitCode := run(nil, repo, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("run(generate) exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if exitCode := run([]string{"--check"}, repo, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("run(check) exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	if exitCode := run([]string{"--unknown"}, repo, &stdout, &stderr); exitCode != 2 {
		t.Errorf("run(unknown) exit code = %d, want 2", exitCode)
	}
}
