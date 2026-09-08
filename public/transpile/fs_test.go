package transpile

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestTranspileFSRejectsDirectoryAccess(t *testing.T) {
	t.Parallel()

	fs := &transpileFS{files: map[string]string{"/src/module.ts": ""}}
	testutil.AssertPanics(t, func() {
		fs.DirectoryExists("/src")
	}, `unexpected directory existence check for "/src"`)
	testutil.AssertPanics(t, func() {
		fs.Realpath("/src/module.ts")
	}, `unexpected realpath request for "/src/module.ts"`)
}
