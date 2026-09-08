package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestOutliningHintSpansForFunction(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `namespace NS[| {
    function f(x: number, y: number)[| {
        return x + y;
    }|]

    function g[|(
        x: number,
        y: number,
    ): number {
        return x + y;
    }|]
}|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOutliningSpans(t)
}
