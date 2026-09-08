package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFindReferencesSeeTagInTs(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function doStuffWithStuff/*1*/(stuff: { quantity: number }) {}

declare const stuff: { quantity: number };
/** @see {doStuffWithStuff} */
if (stuff.quantity) {}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineFindAllReferences(t, "1")
}
