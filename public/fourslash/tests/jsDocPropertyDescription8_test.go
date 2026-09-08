package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocPropertyDescription8(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class SymbolClass {
    /** Something generic */
    static [p: symbol]: any;
}
function symbolClass(e: typeof SymbolClass) {
    console.log(e./*symbolClass*/anything);
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "symbolClass", "any", "")
}
