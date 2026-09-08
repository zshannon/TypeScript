package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestInlayHintsRestParameters2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo(a: unknown, b: unknown, c: unknown) { }
function foo1(...x: [number, number | undefined]) {
    foo(...x, 3);
}
function foo2(...x: []) {
    foo(...x, 1, 2, 3);
}
function foo3(...x: [number, number?]) {
    foo(1, ...x);
}
function foo4(...x: [number, number?]) {
    foo(...x, 3);
}
function foo5(...x: [number, number]) {
    foo(...x, 3);
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineInlayHints(t, nil /*span*/, &lsutil.UserPreferences{InlayHints: lsutil.InlayHintsPreferences{IncludeInlayParameterNameHints: lsutil.IncludeInlayParameterNameHintsAll}})
}
