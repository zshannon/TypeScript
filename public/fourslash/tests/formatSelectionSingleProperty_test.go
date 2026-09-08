package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormatSelectionSingleProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `console.log({
}, {
/*1*/    a: 1,
/*2*/    b: 2
})`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatSelection(t, "1", "2")
	f.VerifyCurrentFileContent(t, `console.log({
}, {
    a: 1,
    b: 2
})`)
}
