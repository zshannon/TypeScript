package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocAugments(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: dummy.js
/**
 * @augments {Thing<string>}
 */
class MyStringThing extends Thing {
    constructor() {
        var x = this.mine;
        x/**/;
    }
}
// @Filename: declarations.d.ts
declare class Thing<T> {
    mine: T;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.VerifyQuickInfoIs(t, "(local var) x: string", "")
}
