package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestFormatSpaceBetweenFunctionAndArrayIndex(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5

function test() {
    return [];
}

test() [0]
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.MarkTestAsStradaServer()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
function test() {
    return [];
}

test()[0]
`)
}
