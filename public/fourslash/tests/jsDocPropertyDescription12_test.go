package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocPropertyDescription12(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type SymbolAlias = {
    /** Something generic */
    [p: symbol]: string;
}
function symbolAlias(e: SymbolAlias) {
    console.log(e./*symbolAlias*/anything);
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "symbolAlias", "any", "")
}
