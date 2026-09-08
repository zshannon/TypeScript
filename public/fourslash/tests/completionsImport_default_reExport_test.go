package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	. "github.com/zshannon/TypeScript/public/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/ls"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestCompletionsImport_default_reExport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
// @module: commonjs
// @allowJs: true
// @Filename: /file1.js
const a = 1;
export {
    a as b
};
export default a;
// @Filename: /file2.js
import * as foo from './file1';
/**/
export default foo.b;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: CompletionGlobalsInJSPlus(
				[]fourslash.CompletionsExpectedItem{
					"foo",
					&lsproto.CompletionItem{
						Label: "a",
						Data: &lsproto.CompletionItemData{
							AutoImport: &lsproto.AutoImportFix{
								ModuleSpecifier: "./file1",
							},
						},
						AdditionalTextEdits: fourslash.AnyTextEdits,
						SortText:            new(string(ls.SortTextAutoImportSuggestions)),
					},
					&lsproto.CompletionItem{
						Label: "b",
						Data: &lsproto.CompletionItemData{
							AutoImport: &lsproto.AutoImportFix{
								ModuleSpecifier: "./file1",
							},
						},
						AdditionalTextEdits: fourslash.AnyTextEdits,
						SortText:            new(string(ls.SortTextAutoImportSuggestions)),
					},
				}, false,
			),
		},
	})
}
