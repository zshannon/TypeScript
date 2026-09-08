package printer

import (
	"github.com/zshannon/TypeScript/public/v7/ast"
	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/tsoptions"
	"github.com/zshannon/TypeScript/public/v7/tspath"
)

// NOTE: EmitHost operations must be thread-safe
type EmitHost interface {
	Options() *core.CompilerOptions
	SourceFiles() []*ast.SourceFile
	UseCaseSensitiveFileNames() bool
	GetCurrentDirectory() string
	CommonSourceDirectory() string
	IsEmitBlocked(file string) bool
	WriteFile(fileName string, text string) error
	GetEmitModuleFormatOfFile(file ast.HasFileName) core.ModuleKind
	GetEmitResolver() EmitResolver
	GetProjectReferenceFromSource(path tspath.Path) *tsoptions.SourceOutputAndProjectReference
	IsSourceFileFromExternalLibrary(file *ast.SourceFile) bool
}
