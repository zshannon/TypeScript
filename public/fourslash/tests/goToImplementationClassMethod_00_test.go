package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestGoToImplementationClassMethod_00(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class Bar {
    [|{|"parts": ["(","method",")"," ","Bar",".","hello","(",")",":"," ","void"], "kind": "method"|}hello|]() {}
}

new Bar().hel/*reference*/lo;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToImplementation(t, "reference")
}
