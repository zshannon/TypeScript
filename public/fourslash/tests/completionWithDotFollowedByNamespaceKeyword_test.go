package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	. "github.com/zshannon/TypeScript/public/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestCompletionWithDotFollowedByNamespaceKeyword(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `namespace A {
    function foo() {
        if (true) {
            B./**/
        namespace B {
            export function baz() { }
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
			Exact: []fourslash.CompletionsExpectedItem{
				&lsproto.CompletionItem{
					Label:  "baz",
					Detail: new("function B.baz(): void"),
				},
			},
		},
	})
}
