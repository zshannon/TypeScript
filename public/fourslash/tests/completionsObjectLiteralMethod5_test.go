package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/ls"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCompletionsObjectLiteralMethod5(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @newline: LF
// @Filename: a.ts
interface Foo {
    method(x?: string): void;
}
const foo: Foo = {
    /*m*/
}`
	f, done := fourslash.NewFourslash(t, fourslash.GetDefaultCapabilitiesWithOptions(&fourslash.ClientCapabilitiesOptions{
		CompletionItem: &lsproto.ClientCompletionItemOptions{
			SnippetSupport:      new(true),
			LabelDetailsSupport: new(true),
		},
	}), content)
	defer done()
	f.VerifyCompletions(t, "m", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:    "method",
					SortText: new(string(ls.ObjectLiteralPropertySortText(ls.SortTextLocationPriority, "method"))),
				},
				&lsproto.CompletionItem{
					Label:      "method",
					InsertText: new("method(x) {\n    $0\n},"),
					SortText:   new(string(ls.SortBelow(ls.ObjectLiteralPropertySortText(ls.SortTextLocationPriority, "method")))),
					Data: &lsproto.CompletionItemData{
						Source: "ObjectLiteralMethodSnippet/",
					},
					InsertTextFormat: new(lsproto.InsertTextFormatSnippet),
					LabelDetails: &lsproto.CompletionItemLabelDetails{
						Detail: new("(x)"),
					},
				},
			},
		},
		UserPreferences: &lsutil.UserPreferences{IncludeCompletionsWithObjectLiteralMethodSnippets: core.TSTrue},
	})
}
