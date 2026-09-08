package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestAutoImportPackageJsonImportsCaseSensitivity(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: node18
// @allowImportingTsExtensions: true
// @Filename: /package.json
{
  "type": "module",
  "imports": {
    "#src/*": "./SRC/*"
  }
}
// @Filename: /src/add.ts
export function add(a: number, b: number) {}
// @Filename: /src/index.ts
add/*imports*/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "imports", []string{"#src/add.ts"}, &lsutil.UserPreferences{ImportModuleSpecifierPreference: "non-relative"})
}
