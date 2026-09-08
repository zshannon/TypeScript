package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	. "github.com/zshannon/TypeScript/public/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestCompletionOfAwaitPromise2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface Foo { foo: string }
async function foo(x: Promise<Foo>) {
   [|x./**/|]
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"then",
				&lsproto.CompletionItem{
					Label:      "foo",
					InsertText: new("(await x).foo"),
					TextEdit: &lsproto.TextEditOrInsertReplaceEdit{
						TextEdit: &lsproto.TextEdit{
							NewText: "(await x).foo",
							Range:   f.Ranges()[0].LSRange,
						},
					},
				},
			},
		},
	})
}
