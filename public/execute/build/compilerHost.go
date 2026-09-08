package build

import (
	"github.com/zshannon/TypeScript/public/v7/ast"
	"github.com/zshannon/TypeScript/public/v7/compiler"
	"github.com/zshannon/TypeScript/public/v7/diagnostics"
	"github.com/zshannon/TypeScript/public/v7/tsoptions"
	"github.com/zshannon/TypeScript/public/v7/tspath"
	"github.com/zshannon/TypeScript/public/v7/vfs"
)

type compilerHost struct {
	host  *host
	trace func(msg *diagnostics.Message, args ...any)
}

var _ compiler.CompilerHost = (*compilerHost)(nil)

func (h *compilerHost) FS() vfs.FS {
	return h.host.FS()
}

func (h *compilerHost) DefaultLibraryPath() string {
	return h.host.DefaultLibraryPath()
}

func (h *compilerHost) GetCurrentDirectory() string {
	return h.host.GetCurrentDirectory()
}

func (h *compilerHost) Trace(msg *diagnostics.Message, args ...any) {
	h.trace(msg, args...)
}

func (h *compilerHost) GetSourceFile(opts ast.SourceFileParseOptions) *ast.SourceFile {
	return h.host.GetSourceFile(opts)
}

func (h *compilerHost) GetResolvedProjectReference(fileName string, path tspath.Path) *tsoptions.ParsedCommandLine {
	return h.host.GetResolvedProjectReference(fileName, path)
}
