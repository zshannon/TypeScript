package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestAutoImportPathsConfigDir(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: tsconfig.json
{
    "compilerOptions": {
        "paths": {
            "@root/*": ["${configDir}/src/*"]
        }
    }
}
// @Filename: src/one.ts
export const one = 1;
// @Filename: src/foo/two.ts
one/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"@root/one"}, nil /*preferences*/)
}
