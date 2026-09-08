package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestRenameStringLiteralOk1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare function f(): '[|foo|]' | 'bar'
class Foo {
    f = f()
}
const d: 'foo' = 'foo'
declare const ff: Foo
ff.f = '[|foo|]'`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineRenameAtRangesWithText(t, nil /*preferences*/, "foo")
}
