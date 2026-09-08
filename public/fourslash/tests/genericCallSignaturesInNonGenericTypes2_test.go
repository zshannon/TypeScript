package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGenericCallSignaturesInNonGenericTypes2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface WrappedArray<T> { }
interface Underscore {
    <T>(list: T[]): WrappedArray<T>;
}
var _: Underscore;
var a: number[];
var /**/b = _(a);  // WrappedArray<any>, should be WrappedArray<number>`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "", "var b: WrappedArray<number>", "")
}
