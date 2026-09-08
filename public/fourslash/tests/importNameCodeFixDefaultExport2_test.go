package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestImportNameCodeFixDefaultExport2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /lib.ts
class Base { }
export default Base;
// @Filename: /test.ts
[|class Derived extends Base { }|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToFile(t, "/test.ts")
	f.VerifyImportFixAtPosition(t, []string{
		`import Base from "./lib";

class Derived extends Base { }`,
	}, nil /*preferences*/)
}
