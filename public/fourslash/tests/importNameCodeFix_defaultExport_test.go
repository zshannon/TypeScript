package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestImportNameCodeFix_defaultExport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @allowJs: true
// @checkJs: true
// @Filename: /a.js
class C {}
export default C;
// @Filename: /b.js
[|C;|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToFile(t, "/b.js")
	f.VerifyImportFixAtPosition(t, []string{
		`import C from "./a";

C;`,
	}, nil /*preferences*/)
}
