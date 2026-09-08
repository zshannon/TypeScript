package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSignatureHelpTaggedTemplatesNegatives4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo(strs, ...rest) {
}

/*1*/fo/*2*/o /*3*/` + "`" + `` + "`" + `/*4*/  /*5*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyNoSignatureHelpForMarkers(t, f.MarkerNames()...)
}
