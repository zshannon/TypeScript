package printer

import (
	"github.com/zshannon/TypeScript/public/ast"
	"github.com/zshannon/TypeScript/public/tspath"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData
}
