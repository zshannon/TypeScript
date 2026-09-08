package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"PUBLIC_MODULE_PATH/ast"
	"PUBLIC_MODULE_PATH/bundled"
	"PUBLIC_MODULE_PATH/compiler"
	"PUBLIC_MODULE_PATH/core"
	"PUBLIC_MODULE_PATH/parser"
	"PUBLIC_MODULE_PATH/tsoptions"
	"PUBLIC_MODULE_PATH/tspath"
	"PUBLIC_MODULE_PATH/vfs/osvfs"
)

const validSource = `
const values: number[] = [1, 2, 3];
const doubled: number[] = values.map(value => value * 2);
`

const invalidSource = `const answer: number = "forty-two";`

func main() {
	expectedVersion := os.Getenv("EXPECTED_TYPESCRIPT_VERSION")
	if expectedVersion == "" {
		fatalf("EXPECTED_TYPESCRIPT_VERSION is required")
	}
	if actualVersion := core.Version(); actualVersion != expectedVersion {
		fatalf("public compiler version is %s, want %s", actualVersion, expectedVersion)
	}

	root, err := os.MkdirTemp("", "typescript-public-consumer-")
	if err != nil {
		fatalf("create TypeScript fixture directory: %v", err)
	}
	defer os.RemoveAll(root)

	parseFile := filepath.Join(root, "parse.ts")
	parsed := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: parseFile,
		Path:     tspath.Path(tspath.NormalizeSlashes(parseFile)),
	}, validSource, core.ScriptKindTS)
	if diagnostics := parsed.Diagnostics(); len(diagnostics) != 0 {
		fatalf("public parser returned %d diagnostics for valid TypeScript", len(diagnostics))
	}
	if len(parsed.Statements.Nodes) != 2 {
		fatalf("public parser returned %d statements, want 2", len(parsed.Statements.Nodes))
	}

	checkProgram(root, "valid", validSource, false)
	checkProgram(root, "invalid", invalidSource, true)
	fmt.Println("PASS: unrelated module parsed and typechecked TypeScript through public packages")
}

func checkProgram(root string, name string, source string, wantTypeError bool) {
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fatalf("create %s program directory: %v", name, err)
	}
	sourcePath := filepath.Join(dir, "index.ts")
	configPath := filepath.Join(dir, "tsconfig.json")
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		fatalf("write %s TypeScript source: %v", name, err)
	}
	config := `{"compilerOptions":{"strict":true,"target":"ES2022","noEmit":true},"files":["index.ts"]}`
	if err := os.WriteFile(configPath, []byte(config), 0o644); err != nil {
		fatalf("write %s tsconfig: %v", name, err)
	}

	currentDirectory := tspath.NormalizeSlashes(dir)
	fs := bundled.WrapFS(osvfs.FS())
	host := compiler.NewCompilerHost(currentDirectory, fs, bundled.LibPath(), nil, nil)
	parsed, configDiagnostics := tsoptions.GetParsedCommandLineOfConfigFile(
		tspath.NormalizeSlashes(configPath),
		&core.CompilerOptions{},
		nil,
		host,
		nil,
	)
	if len(configDiagnostics) != 0 {
		fatalf("%s config returned %d diagnostics", name, len(configDiagnostics))
	}

	program := compiler.NewProgram(compiler.ProgramOptions{Config: parsed, Host: host})
	if !loadsStandardLibrary(program) {
		fatalf("%s program did not resolve an embedded standard library", name)
	}
	ctx := context.Background()
	if diagnostics := program.GetSyntacticDiagnostics(ctx, nil); len(diagnostics) != 0 {
		fatalf("%s program returned %d syntactic diagnostics", name, len(diagnostics))
	}
	if diagnostics := program.GetProgramDiagnostics(); len(diagnostics) != 0 {
		fatalf("%s program returned %d program diagnostics", name, len(diagnostics))
	}
	program.BindSourceFiles()
	if diagnostics := program.GetGlobalDiagnostics(ctx); len(diagnostics) != 0 {
		fatalf("%s program returned %d global diagnostics", name, len(diagnostics))
	}
	semanticDiagnostics := program.GetSemanticDiagnostics(ctx, nil)
	if !wantTypeError {
		if len(semanticDiagnostics) != 0 {
			fatalf("valid program returned %d semantic diagnostics", len(semanticDiagnostics))
		}
		return
	}
	if len(semanticDiagnostics) != 1 || semanticDiagnostics[0].Code() != 2322 {
		codes := make([]string, 0, len(semanticDiagnostics))
		for _, diagnostic := range semanticDiagnostics {
			codes = append(codes, fmt.Sprint(diagnostic.Code()))
		}
		fatalf("invalid program returned diagnostic codes [%s], want [2322]", strings.Join(codes, ", "))
	}
}

func loadsStandardLibrary(program *compiler.Program) bool {
	for _, sourceFile := range program.GetSourceFiles() {
		if strings.HasSuffix(sourceFile.FileName(), "/lib.es5.d.ts") {
			return true
		}
	}
	return false
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", args...)
	os.Exit(1)
}
