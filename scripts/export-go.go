package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	publicModuleBase = "github.com/zshannon/TypeScript/public"
)

var errOutputOutOfDate = errors.New("generated public module is out of date")

const publicReadmeTemplate = `# TypeScript compiler packages for Go

This module is generated from every package under ` + "`tsc/internal`" + ` in the
[TypeScript Go repository](https://github.com/zshannon/TypeScript). It publishes
the whole compiler package surface with normal Go visibility: exported identifiers
remain exported, while unexported identifiers remain private to their packages.

For example:

` + "```go" + `
import (
	"%[1]s/ast"
	"%[1]s/parser"
)
` + "```" + `

Nested source directories named ` + "`internal`" + ` are generated as ` + "`internals`" + ` so
packages such as ` + "`vfs/internal`" + ` can be imported by external modules.

## Compatibility

This broad API follows upstream compiler implementation packages and may change at
any time. Go module versions align exactly with upstream TypeScript versions. For
example, ` + "`github.com/zshannon/TypeScript/public/v7@v7.0.2`" + ` corresponds to TypeScript 7.0.2.

%[2]s

## Regeneration

From the repository root:

` + "```sh" + `
go run ./scripts/export-go.go
go run ./scripts/export-go.go --check
go run ./scripts/export-go.go --version
` + "```" + `

The generator copies compiler packages and their required source assets. It does not
copy the repository-level ` + "`tsc/testdata`" + ` corpus; upstream test harness packages that
load those fixtures at runtime require a full TypeScript repository checkout.
`

func main() {
	os.Exit(run(os.Args[1:], ".", os.Stdout, os.Stderr))
}

func run(args []string, repoRoot string, stdout, stderr io.Writer) int {
	mode := "generate"
	switch {
	case len(args) == 0:
	case len(args) == 1 && args[0] == "--check":
		mode = "check"
	case len(args) == 1 && args[0] == "--version":
		mode = "version"
	default:
		fmt.Fprintln(stderr, "usage: go run ./scripts/export-go.go [--check|--version]")
		return 2
	}

	root, err := findRepoRoot(repoRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if mode == "version" {
		version, err := readCompilerVersion(root)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(stdout, "v"+version.source)
		return 0
	}
	if err := exportModule(root, mode == "check"); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if mode == "check" {
		fmt.Fprintln(stdout, "public module is up to date")
	} else {
		fmt.Fprintln(stdout, "generated public module")
	}
	return 0
}

func exportModule(repoRoot string, check bool) error {
	expected, err := buildOutput(repoRoot)
	if err != nil {
		return err
	}
	destination := filepath.Join(repoRoot, "public")
	if check {
		differences, err := compareOutput(destination, expected)
		if err != nil {
			return err
		}
		if len(differences) != 0 {
			return fmt.Errorf("%w:\n  %s", errOutputOutOfDate, strings.Join(differences, "\n  "))
		}
		return nil
	}
	return writeOutput(destination, expected)
}

type outputFile struct {
	data []byte
	mode fs.FileMode
}

type compilerVersion struct {
	source       string
	major        int
	publicModule string
}

func buildOutput(repoRoot string) (map[string]outputFile, error) {
	version, err := readCompilerVersion(repoRoot)
	if err != nil {
		return nil, err
	}
	goMod, err := os.ReadFile(filepath.Join(repoRoot, "tsc", "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("read tsc/go.mod: %w", err)
	}
	sourceModule, err := readModuleDirective(goMod)
	if err != nil {
		return nil, err
	}
	sourceRoot := filepath.Join(repoRoot, "tsc", "internal")
	if info, err := os.Stat(sourceRoot); err != nil || !info.IsDir() {
		if err == nil {
			err = fmt.Errorf("not a directory")
		}
		return nil, fmt.Errorf("compiler package source %q: %w", sourceRoot, err)
	}

	result := map[string]outputFile{}
	err = filepath.WalkDir(sourceRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("compiler package source %q is not a regular file", path)
		}
		rel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		outputPath := mapInternalSegments(rel)
		if _, exists := result[outputPath]; exists {
			return fmt.Errorf("compiler package paths collide at generated path %q", filepath.ToSlash(outputPath))
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".go") {
			content, err = rewriteGoImports(path, content, sourceModule, version.publicModule)
			if err != nil {
				return err
			}
		}
		result[outputPath] = outputFile{data: content, mode: info.Mode().Perm()}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect compiler packages: %w", err)
	}

	goMod, err = rewriteModuleDirective(goMod, version.publicModule)
	if err != nil {
		return nil, err
	}
	result["go.mod"] = outputFile{data: goMod, mode: 0o644}

	for _, file := range []struct {
		source string
		output string
	}{
		{source: filepath.Join(repoRoot, "tsc", "go.sum"), output: "go.sum"},
		{source: filepath.Join(repoRoot, "tsc", "LICENSE"), output: "LICENSE"},
	} {
		content, err := os.ReadFile(file.source)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file.source, err)
		}
		result[file.output] = outputFile{data: content, mode: 0o644}
	}
	notice, err := readFirstExisting(
		filepath.Join(repoRoot, "tsc", "NOTICE.txt"),
		filepath.Join(repoRoot, "NOTICE.txt"),
	)
	if err != nil {
		return nil, err
	}
	result["NOTICE"] = outputFile{data: notice, mode: 0o644}
	result["README.md"] = outputFile{data: []byte(publicReadme(version)), mode: 0o644}
	return result, nil
}

func readFirstExisting(paths ...string) ([]byte, error) {
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err == nil {
			return content, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
	}
	return nil, fmt.Errorf("none of the required files exist: %s", strings.Join(paths, ", "))
}

func readCompilerVersion(repoRoot string) (compilerVersion, error) {
	versionPath := filepath.Join(repoRoot, "tsc", "internal", "core", "version.go")
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, versionPath, nil, 0)
	if err != nil {
		return compilerVersion{}, fmt.Errorf("read compiler version from %s: %w", versionPath, err)
	}
	var source string
	for _, declaration := range parsed.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.VAR {
			continue
		}
		for _, spec := range general.Specs {
			value, ok := spec.(*ast.ValueSpec)
			if !ok || len(value.Names) != 1 || value.Names[0].Name != "version" || len(value.Values) != 1 {
				continue
			}
			literal, ok := value.Values[0].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				return compilerVersion{}, fmt.Errorf("compiler version in %s is not a string literal", versionPath)
			}
			source, err = strconv.Unquote(literal.Value)
			if err != nil {
				return compilerVersion{}, fmt.Errorf("parse compiler version in %s: %w", versionPath, err)
			}
		}
	}
	if source == "" {
		return compilerVersion{}, fmt.Errorf("compiler version variable not found in %s", versionPath)
	}
	majorText, _, found := strings.Cut(source, ".")
	if !found || strings.Count(source, ".") < 2 {
		return compilerVersion{}, fmt.Errorf("compiler version %q does not contain major, minor, and patch components", source)
	}
	major, err := strconv.Atoi(majorText)
	if err != nil || major < 0 {
		return compilerVersion{}, fmt.Errorf("compiler version %q has an invalid major version", source)
	}
	modulePath := publicModuleBase
	if major >= 2 {
		modulePath += "/v" + majorText
	}
	return compilerVersion{source: source, major: major, publicModule: modulePath}, nil
}

func publicReadme(version compilerVersion) string {
	current := fmt.Sprintf("The current generated source is `%s`, and its corresponding Go module version is `v%s`.", version.source, version.source)
	if strings.HasSuffix(version.source, "-dev") {
		current = fmt.Sprintf("The current generated source is `%s`, a development version; it is not the stable 7.0.2 release.", version.source)
	}
	return fmt.Sprintf(publicReadmeTemplate, version.publicModule, current)
}

func mapInternalSegments(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for index, part := range parts {
		if part == "internal" {
			parts[index] = "internals"
		}
	}
	return filepath.FromSlash(strings.Join(parts, "/"))
}

func rewriteGoImports(filename string, content []byte, sourceModule string, publicModule string) ([]byte, error) {
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filename, content, 0)
	if err != nil {
		return nil, fmt.Errorf("parse Go source %s: %w", filename, err)
	}
	type edit struct {
		start int
		end   int
		value string
	}
	var edits []edit
	ast.Inspect(parsed, func(node ast.Node) bool {
		literal, ok := node.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}
		value := rewriteCompilerPathOccurrences(literal.Value, sourceModule, publicModule)
		if value == literal.Value {
			return true
		}
		edits = append(edits, edit{
			start: fileSet.Position(literal.Pos()).Offset,
			end:   fileSet.Position(literal.End()).Offset,
			value: value,
		})
		return true
	})
	for index := len(edits) - 1; index >= 0; index-- {
		item := edits[index]
		content = bytes.Join([][]byte{content[:item.start], []byte(item.value), content[item.end:]}, nil)
	}
	return content, nil
}

func rewriteCompilerPathOccurrences(literal string, sourceModule string, publicModule string) string {
	oldPrefix := sourceModule + "/internal"
	remaining := literal
	var rewritten strings.Builder
	for {
		index := strings.Index(remaining, oldPrefix)
		if index < 0 {
			rewritten.WriteString(remaining)
			return rewritten.String()
		}
		rewritten.WriteString(remaining[:index])
		end := index + len(oldPrefix)
		for end < len(remaining) && isImportPathByte(remaining[end]) {
			end++
		}
		path := remaining[index:end]
		tail := strings.TrimPrefix(path, oldPrefix)
		tail = strings.TrimPrefix(tail, "/")
		rewritten.WriteString(publicModule)
		if tail != "" {
			rewritten.WriteByte('/')
			rewritten.WriteString(filepath.ToSlash(mapInternalSegments(filepath.FromSlash(tail))))
		}
		remaining = remaining[end:]
	}
}

func readModuleDirective(content []byte) (string, error) {
	for _, line := range bytes.Split(content, []byte("\n")) {
		fields := strings.Fields(string(line))
		if len(fields) == 0 || fields[0] != "module" {
			continue
		}
		if len(fields) != 2 {
			return "", errors.New("tsc/go.mod has an invalid module directive")
		}
		return fields[1], nil
	}
	return "", errors.New("tsc/go.mod has no module directive")
}

func isImportPathByte(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9' ||
		strings.ContainsRune("/._~-+", rune(value))
}

func rewriteModuleDirective(content []byte, publicModule string) ([]byte, error) {
	lines := bytes.SplitAfter(content, []byte("\n"))
	for index, line := range lines {
		trimmed := strings.TrimSpace(string(line))
		if !strings.HasPrefix(trimmed, "module ") {
			continue
		}
		lineEnding := ""
		if bytes.HasSuffix(line, []byte("\r\n")) {
			lineEnding = "\r\n"
		} else if bytes.HasSuffix(line, []byte("\n")) {
			lineEnding = "\n"
		}
		lines[index] = []byte("module " + publicModule + lineEnding)
		return bytes.Join(lines, nil), nil
	}
	return nil, errors.New("tsc/go.mod has no module directive")
}

func compareOutput(destination string, expected map[string]outputFile) ([]string, error) {
	actualFiles := map[string]struct{}{}
	expectedDirs := map[string]struct{}{filepath.Clean(destination): {}}
	for rel := range expected {
		for directory := filepath.Dir(filepath.Join(destination, rel)); directory != filepath.Clean(destination); directory = filepath.Dir(directory) {
			expectedDirs[directory] = struct{}{}
		}
	}

	var differences []string
	err := filepath.WalkDir(destination, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) && path == destination {
				differences = append(differences, "missing public output directory")
				return fs.SkipAll
			}
			return walkErr
		}
		if path == destination {
			return nil
		}
		if entry.IsDir() {
			if _, exists := expectedDirs[path]; !exists {
				rel, _ := filepath.Rel(destination, path)
				differences = append(differences, "unexpected directory "+filepath.ToSlash(rel))
				return fs.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(destination, path)
		if err != nil {
			return err
		}
		actualFiles[rel] = struct{}{}
		want, exists := expected[rel]
		if !exists {
			differences = append(differences, "unexpected file "+filepath.ToSlash(rel))
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			differences = append(differences, "non-regular file "+filepath.ToSlash(rel))
			return nil
		}
		got, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Equal(got, want.data) {
			differences = append(differences, "changed file "+filepath.ToSlash(rel))
		}
		if info.Mode().Perm() != want.mode.Perm() {
			differences = append(differences, "changed mode "+filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("inspect generated public module: %w", err)
	}
	for rel := range expected {
		if _, exists := actualFiles[rel]; !exists {
			differences = append(differences, "missing file "+filepath.ToSlash(rel))
		}
	}
	sort.Strings(differences)
	return differences, nil
}

func writeOutput(destination string, expected map[string]outputFile) (returnErr error) {
	parent := filepath.Dir(destination)
	if filepath.Base(destination) != "public" {
		return fmt.Errorf("refusing to replace unexpected output path %q", destination)
	}
	staging, err := os.MkdirTemp(parent, ".public-export-")
	if err != nil {
		return fmt.Errorf("create public output staging directory: %w", err)
	}
	defer func() {
		if cleanupErr := os.RemoveAll(staging); returnErr == nil && cleanupErr != nil {
			returnErr = fmt.Errorf("remove public output staging directory: %w", cleanupErr)
		}
	}()

	paths := make([]string, 0, len(expected))
	for path := range expected {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, rel := range paths {
		file := expected[rel]
		path := filepath.Join(staging, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create generated directory: %w", err)
		}
		if err := os.WriteFile(path, file.data, file.mode); err != nil {
			return fmt.Errorf("write generated file %s: %w", filepath.ToSlash(rel), err)
		}
		if err := os.Chmod(path, file.mode); err != nil {
			return fmt.Errorf("set generated file mode %s: %w", filepath.ToSlash(rel), err)
		}
	}

	backup := ""
	if _, err := os.Lstat(destination); err == nil {
		backup, err = os.MkdirTemp(parent, ".public-backup-")
		if err != nil {
			return fmt.Errorf("reserve public output backup path: %w", err)
		}
		if err := os.Remove(backup); err != nil {
			return fmt.Errorf("prepare public output backup path: %w", err)
		}
		if err := os.Rename(destination, backup); err != nil {
			return fmt.Errorf("back up existing public output: %w", err)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("inspect existing public output: %w", err)
	}

	if err := os.Rename(staging, destination); err != nil {
		if backup != "" {
			_ = os.Rename(backup, destination)
		}
		return fmt.Errorf("install generated public output: %w", err)
	}
	staging = ""
	if backup != "" {
		if err := os.RemoveAll(backup); err != nil {
			return fmt.Errorf("remove previous public output: %w", err)
		}
	}
	return nil
}

func findRepoRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if info, err := os.Stat(filepath.Join(current, "tsc", "internal")); err == nil && info.IsDir() {
			if info, err := os.Stat(filepath.Join(current, "tsc", "go.mod")); err == nil && !info.IsDir() {
				return current, nil
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("could not find repository root from %q", start)
		}
		current = parent
	}
}
