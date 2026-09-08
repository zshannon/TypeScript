package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestEmptyArrayInference(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strict: false
var x/*1*/x = true ? [1] : [undefined]; 
var y/*2*/y = true ? [1] : [];`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "1", "var xx: number[]", "")
	f.VerifyQuickInfoAt(t, "2", "var yy: number[]", "")
}
