package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestAutoImportSpecifierExcludeRegexes2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "module": "preserve",
        "paths": {
            "@app/*": ["./src/*"]
        }
    }
}
// @Filename: /src/utils.ts
export function add(a: number, b: number) {}
// @Filename: /src/index.ts
add/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"./utils"}, nil /*preferences*/)
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"@app/utils"}, &lsutil.UserPreferences{AutoImportSpecifierExcludeRegexes: []string{"^\\./"}})
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"@app/utils"}, &lsutil.UserPreferences{ImportModuleSpecifierPreference: "non-relative"})
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"./utils"}, &lsutil.UserPreferences{ImportModuleSpecifierPreference: "non-relative", AutoImportSpecifierExcludeRegexes: []string{"^@app/"}})
	f.VerifyImportFixModuleSpecifiers(t, "", []string{}, &lsutil.UserPreferences{AutoImportSpecifierExcludeRegexes: []string{"utils"}})
}
