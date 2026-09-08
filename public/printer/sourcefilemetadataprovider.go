package printer

import (
	"github.com/zshannon/TypeScript/public/v7/ast"
	"github.com/zshannon/TypeScript/public/v7/tspath"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData
}
