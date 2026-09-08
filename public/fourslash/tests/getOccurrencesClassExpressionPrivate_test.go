package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGetOccurrencesClassExpressionPrivate(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let A = class Foo {
    [|private|] foo;
    [|private|] private;
    constructor([|private|] y: string, public x: string) {
    }
    [|private|] method() { }
    public method2() { }
    [|private|] static static() { }
}

let B = class D {
    constructor(private x: number) {
    }
    private test() {}
    public test2() {}
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineDocumentHighlights(t, nil /*preferences*/, ToAny(f.Ranges())...)
}
