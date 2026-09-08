package sourcemap

import "github.com/zshannon/TypeScript/public/v7/core"

type Source interface {
	Text() string
	FileName() string
	ECMALineMap() []core.TextPos
}
