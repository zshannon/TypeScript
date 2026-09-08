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
	"strings"
)

const (
	sourceModule = "github.com/microsoft/TypeScript/tsc"
	publicModule = "github.com/zshannon/TypeScript/public"
)

var errOutputOutOfDate = errors.New("generated public module is out of date")

const publicReadme = `# TypeScript compiler packages for Go

This module is generated from every package under ` + "`tsc/internal`" + ` in the
[TypeScript Go repository](https://github.com/zshannon/TypeScript). It publishes
the whole compiler package surface with normal Go visibility: exported identifiers
remain exported, while unexported identifiers remain private to their packages.

For example:

` + "```go" + `
import (
	"github.com/zshannon/TypeScript/public/ast"
	"github.com/zshannon/TypeScript/public/parser"
)
` + "```" + `

Nested source directories named ` + "`internal`" + ` are generated as ` + "`internals`" + ` so
packages such as ` + "`vfs/internal`" + ` can be imported by external modules.

## Compatibility

This broad API follows upstream compiler implementation packages and may change at
any time. Consumers should pin a commit using its Go pseudo-version instead of
assuming semantic API stability. The module uses ordinary Go module resolution
from repository commits and has no separate release process.

## Regeneration

From the repository root:

` + "```sh" + `
go run ./scripts/export-go.go
go run ./scripts/export-go.go --check
` + "```" + `

The generator copies compiler packages and their required source assets. It does not
copy the repository-level ` + "`tsc/testdata`" + ` corpus; upstream test harness packages that
load those fixtures at runtime require a full TypeScript repository checkout.
`

func main() {
	os.Exit(run(os.Args[1:], ".", os.Stdout, os.Stderr))
}

func run(args []string, repoRoot string, stdout, stderr io.Writer) int {
	check := false
	switch {
	case len(args) == 0:
	case len(args) == 1 && args[0] == "--check":
		check = true
	default:
		fmt.Fprintln(stderr, "usage: go run ./scripts/export-go.go [--check]")
		return 2
	}

	root, err := findRepoRoot(repoRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := exportModule(root, check); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if check {
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

func buildOutput(repoRoot string) (map[string]outputFile, error) {
	sourceRoot := filepath.Join(repoRoot, "tsc", "internal")
	if info, err := os.Stat(sourceRoot); err != nil || !info.IsDir() {
		if err == nil {
			err = fmt.Errorf("not a directory")
		}
		return nil, fmt.Errorf("compiler package source %q: %w", sourceRoot, err)
	}

	result := map[string]outputFile{}
	err := filepath.WalkDir(sourceRoot, func(path string, entry fs.DirEntry, walkErr error) error {
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
			content, err = rewriteGoImports(path, content)
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

	goMod, err := os.ReadFile(filepath.Join(repoRoot, "tsc", "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("read tsc/go.mod: %w", err)
	}
	goMod, err = rewriteModuleDirective(goMod)
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
		{source: filepath.Join(repoRoot, "NOTICE.txt"), output: "NOTICE"},
	} {
		content, err := os.ReadFile(file.source)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file.source, err)
		}
		result[file.output] = outputFile{data: content, mode: 0o644}
	}
	result["README.md"] = outputFile{data: []byte(publicReadme), mode: 0o644}
	return result, nil
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

func rewriteGoImports(filename string, content []byte) ([]byte, error) {
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
		value := rewriteCompilerPathOccurrences(literal.Value)
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

func rewriteCompilerPathOccurrences(literal string) string {
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

func isImportPathByte(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9' ||
		strings.ContainsRune("/._~-+", rune(value))
}

func rewriteModuleDirective(content []byte) ([]byte, error) {
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
