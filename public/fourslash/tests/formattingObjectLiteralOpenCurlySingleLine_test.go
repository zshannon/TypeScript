package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormattingObjectLiteralOpenCurlySingleLine(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
let obj1 =
{ x: 10 };

let obj2 =
    // leading trivia
{ y: 10 };
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `
let obj1 =
    { x: 10 };

let obj2 =
    // leading trivia
    { y: 10 };
`)
}
