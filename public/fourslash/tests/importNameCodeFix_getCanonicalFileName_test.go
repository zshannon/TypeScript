package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestImportNameCodeFix_getCanonicalFileName(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /howNow/node_modules/brownCow/index.d.ts
export const foo: number;
// @Filename: /howNow/a.ts
foo;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToFile(t, "/howNow/a.ts")
	f.VerifyImportFixAtPosition(t, []string{
		`import { foo } from "brownCow";

foo;`,
	}, nil /*preferences*/)
}
