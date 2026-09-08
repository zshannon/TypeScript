package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestRenameJsPrototypeProperty02(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: a.js
function bar() {
}
[|bar.prototype.[|{| "contextRangeIndex": 0 |}x|] = 10;|]
var t = new bar();
[|t.[|{| "contextRangeIndex": 2 |}x|] = 11;|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineRenameAtRangesWithText(t, nil /*preferences*/, "x")
}
