package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormattingElseInsideAFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var x = function() {
    if (true) {
    /*1*/} else {/*2*/
}

// newline at the end of the file`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "2")
	f.InsertLine(t, "")
	f.GoToMarker(t, "1")
	f.VerifyCurrentLineContent(t, `    } else {`)
}
