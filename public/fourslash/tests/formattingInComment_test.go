package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormattingInComment(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class A {
foo(              ); // /*1*/
}
function foo() {       var x;       } // /*2*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, ";")
	f.VerifyCurrentLineContent(t, `foo(              ); // ;`)
	f.GoToMarker(t, "2")
	f.Insert(t, "}")
	f.VerifyCurrentLineContent(t, `function foo() {       var x;       } // }`)
}
