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

func TestCompletionsOverridingMethod18(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: a.ts
// @newline: LF
declare function decorator(...args: any[]): any;
class DecoratorBase {
    protected foo(a: string): string;
    protected foo(a: number): number;
    protected foo(a: any): any {
        return a;
    }
}
class DecoratorSub extends DecoratorBase {
    @decorator protected /**/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &[]string{},
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:               "foo",
					InsertText:          new("protected foo(a: string): string;\nprotected foo(a: number): number;\n@decorator\nprotected foo(a: any) {\n}"),
					FilterText:          new("foo"),
					SortText:            new(string(ls.SortTextLocationPriority)),
					AdditionalTextEdits: fourslash.AnyTextEdits,
					Data: &lsproto.CompletionItemData{
						Source: "ClassMemberSnippet/",
					},
				},
			},
		},
		UserPreferences: &lsutil.UserPreferences{IncludeCompletionsWithClassMemberSnippets: core.TSTrue},
	})
	f.VerifyApplyCodeActionFromCompletion(t, new(""), &fourslash.ApplyCodeActionFromCompletionOptions{
		UserPreferences: &lsutil.UserPreferences{IncludeCompletionsWithClassMemberSnippets: core.TSTrue},
		Name:            "foo",
		Source:          "ClassMemberSnippet/",
		Description:     "Update modifiers of 'foo'",
		NewFileContent: new(`declare function decorator(...args: any[]): any;
class DecoratorBase {
    protected foo(a: string): string;
    protected foo(a: number): number;
    protected foo(a: any): any {
        return a;
    }
}
class DecoratorSub extends DecoratorBase {
    
}`),
	})
}
