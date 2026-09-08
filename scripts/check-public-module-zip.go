package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/module"
	modzip "golang.org/x/mod/zip"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "usage: check-public-module-zip <public-module-directory> <module-path> <source-version> <extraction-directory>")
		os.Exit(2)
	}
	version := module.Version{Path: os.Args[2], Version: os.Args[3]}
	if err := module.Check(version.Path, version.Version); err != nil {
		fatalf("invalid public module path or source version: %v", err)
	}
	moduleDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		fatalf("resolve public module directory: %v", err)
	}
	if _, err := modzip.CheckDir(moduleDir); err != nil {
		fatalf("public module cannot be packaged: %v", err)
	}

	testRoot, err := os.MkdirTemp("", "typescript-public-module-zip-")
	if err != nil {
		fatalf("create module zip directory: %v", err)
	}
	defer os.RemoveAll(testRoot)

	zipPath := filepath.Join(testRoot, "public.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		fatalf("create module zip: %v", err)
	}
	if err := modzip.CreateFromDir(zipFile, version, moduleDir); err != nil {
		zipFile.Close()
		fatalf("create module zip from public directory: %v", err)
	}
	if err := zipFile.Close(); err != nil {
		fatalf("close module zip: %v", err)
	}

	extractedDir := os.Args[4]
	if err := modzip.Unzip(extractedDir, version, zipPath); err != nil {
		fatalf("extract generated module zip: %v", err)
	}
	if _, err := os.Stat(filepath.Join(extractedDir, "go.mod")); err != nil {
		fatalf("extracted module is missing go.mod: %v", err)
	}
	fmt.Println("PASS: public module directory forms a valid extractable Go module zip")
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", args...)
	os.Exit(1)
}
