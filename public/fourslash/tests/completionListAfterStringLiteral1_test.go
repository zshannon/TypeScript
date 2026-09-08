package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/ls"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCompletionListAfterStringLiteral1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
"a"./**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Unsorted: []fourslash.CompletionsExpectedItem{
				"toString",
				"charAt",
				"charCodeAt",
				"concat",
				"indexOf",
				"lastIndexOf",
				"localeCompare",
				"match",
				"replace",
				"search",
				"slice",
				"split",
				"substring",
				"toLowerCase",
				"toLocaleLowerCase",
				"toUpperCase",
				"toLocaleUpperCase",
				"trim",
				"length",
				&lsproto.CompletionItem{
					Label:    "substr",
					SortText: new(string(ls.DeprecateSortText(ls.SortTextLocationPriority))),
					Tags:     &[]lsproto.CompletionItemTag{lsproto.CompletionItemTagDeprecated},
				},
				"valueOf",
			},
		},
	})
}
