package sourcemap

import "github.com/zshannon/TypeScript/public/core"

type Source interface {
	Text() string
	FileName() string
	ECMALineMap() []core.TextPos
}
