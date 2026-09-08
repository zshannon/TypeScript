package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocPropertyDescription11(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type AliasExample = {
    /** Something generic */
    [p: string]: string;
    /** Something else */
    [key: ` + "`" + `any${string}` + "`" + `]: string;
}
function aliasExample(e: AliasExample) {
    console.log(e./*alias*/anything);
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "alias", "(index) AliasExample[string | `any${string}`]: string", "Something generic\nSomething else")
}
