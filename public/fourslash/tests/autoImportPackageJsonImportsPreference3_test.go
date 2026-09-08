package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestAutoImportPackageJsonImportsPreference3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: node18
// @Filename: /package.json
{
  "imports": {
    "#*": "./src/*.ts"
  }
}
// @Filename: /src/a/b/c/something.ts
export function something(name: string): any;
// @Filename: /src/a/b/c/d.ts
something/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"#a/b/c/something"}, &lsutil.UserPreferences{ImportModuleSpecifierPreference: "non-relative"})
}
