package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestQuickInfoJSDocParamWithTrailingAtBeforeCommentEnd(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: /a.js
/** @param {string} x trailing @/*at*/*/
function /*fn*/foo(/*x*/x) {}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "fn", "function foo(x: string): void", "\n\n*@param* `x` — trailing @")
	f.VerifyQuickInfoAt(t, "x", "(parameter) x: string", "trailing @")
	f.VerifyCompletions(t, "at", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label: "param",
					Kind:  new(lsproto.CompletionItemKindKeyword),
				},
			},
		},
	})
}
