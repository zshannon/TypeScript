package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	. "github.com/zshannon/TypeScript/public/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestImportCompletionsPackageJsonImportsPatternRootWildcard(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: nodenext
// @Filename: /package.json
{
  "imports": {
    "#/*": "./src/*"
  }
}
// @Filename: /src/something.ts
export function something(name: string): any;
// @Filename: /src/features/bar.ts
export function bar(): any;
// @Filename: /a.ts
import {} from "#//*1*/";`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, []string{"1"}, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"something.js",
				"features",
			},
		},
	})
}
