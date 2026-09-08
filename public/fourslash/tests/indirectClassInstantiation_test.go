package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	. "github.com/zshannon/TypeScript/public/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/ls"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestIndirectClassInstantiation(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: something.js
function TestObj(){
    this.property = "value";
}
var constructor = TestObj;
var instance = new constructor();
instance./*a*/
var class2 = function() { };
class2.prototype.blah = function() { };
var inst2 = new class2();
inst2.blah/*b*/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "a")
	f.VerifyCompletions(t, nil, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: []fourslash.CompletionsExpectedItem{
				"property",
				&lsproto.CompletionItem{
					Label:    "blah",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "class2",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "constructor",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "inst2",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "instance",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "prototype",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
				&lsproto.CompletionItem{
					Label:    "TestObj",
					SortText: new(string(ls.SortTextJavascriptIdentifiers)),
				},
			},
		},
	})
	f.Backspace(t, 1)
	f.GoToMarker(t, "b")
	f.VerifyQuickInfoIs(t, "(method) class2.blah(): void", "")
}
