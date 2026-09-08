package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSemicolonFormatting(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/**/function of1 (b:{r:{c:number`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToEOF(t)
	f.Insert(t, ";")
	f.VerifyCurrentLineContent(t, `function of1(b: { r: { c: number;`)
}
