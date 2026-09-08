package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestAddDuplicateSetter(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class C {
    set foo(value) { }
    /**/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Insert(t, "set foo(value) { }")
}
