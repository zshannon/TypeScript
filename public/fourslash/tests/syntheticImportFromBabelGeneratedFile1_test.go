package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestSyntheticImportFromBabelGeneratedFile1(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @allowSyntheticDefaultImports: true
// @Filename: /a.js
exports.__esModule = true;
exports.default = f;
/**
 * Run this function
 * @param {string} t
 */
function f(t) {}
// @Filename: /b.js
import f from "./a"
/**/f`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "", "(alias) function f(t: string): void\nimport f", "Run this function")
}
