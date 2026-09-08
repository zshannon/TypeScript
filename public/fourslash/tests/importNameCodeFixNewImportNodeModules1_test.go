package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestImportNameCodeFixNewImportNodeModules1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `[|f1/*0*/();|]
// @Filename: ../package.json
{ "dependencies": { "fake-module": "latest" } }
// @Filename: ../node_modules/fake-module/nested.ts
export var v1 = 5;
export function f1();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.Configure(t, lsutil.UserPreferences{AutoImportEntrypointDirectorySearch: core.TSTrue})
	f.VerifyImportFixAtPosition(t, []string{
		`import { f1 } from "fake-module/nested";

f1();`,
	}, nil /*preferences*/)
}
