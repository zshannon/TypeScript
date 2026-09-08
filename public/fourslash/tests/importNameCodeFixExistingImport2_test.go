package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestImportNameCodeFixExistingImport2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import * as ns from "./module";
// Comment
f1/*0*/();
// @Filename: module.ts
 export function f1() {}
 export var v1 = 5;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixAtPosition(t, []string{
		`import * as ns from "./module";
// Comment
ns.f1();`,
		`import * as ns from "./module";
import { f1 } from "./module";
// Comment
f1();`,
	}, nil /*preferences*/)
}
