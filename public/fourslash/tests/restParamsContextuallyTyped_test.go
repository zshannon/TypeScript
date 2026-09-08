package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestRestParamsContextuallyTyped(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var foo: Function = function (/*1*/a, /*2*/b, /*3*/c) { };`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "1", "(parameter) a: any", "")
	f.VerifyQuickInfoAt(t, "2", "(parameter) b: any", "")
	f.VerifyQuickInfoAt(t, "3", "(parameter) c: any", "")
}
