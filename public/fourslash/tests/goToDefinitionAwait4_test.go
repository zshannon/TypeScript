package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestGoToDefinitionAwait4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `async function outerAsyncFun() {
    let /*end*/af = async () => {
      [|/*start*/await|] Promise.resolve(0);
    }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToDefinition(t, true, "start")
}
