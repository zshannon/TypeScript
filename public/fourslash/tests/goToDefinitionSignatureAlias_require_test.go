package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGoToDefinitionSignatureAlias_require(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: /a.js
module.exports = function /*f*/f() {}
// @Filename: /b.js
const f = require("./a");
[|/*use*/f|]();
// @Filename: /bar.ts
import f = require("./a");
[|/*useTs*/f|]();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToDefinition(t, true, "use", "useTs")
}
