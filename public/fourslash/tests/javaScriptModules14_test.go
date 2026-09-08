package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/ls"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJavaScriptModules14(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: myMod.js
if (true) {
    exports.b = true;
} else {
    exports.n = 3;
}
function fn() {
    exports.s = 'foo';
}
var invisible = true;
// @Filename: isGlobal.js
var y = 10;
// @Filename: consumer.js
var x = require('myMod');
/**/;`
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
				&lsproto.CompletionItem{
					Label:    "y",
					SortText: new(string(ls.SortTextGlobalsOrKeywords)),
				},
			},
			Excludes: []string{
				"invisible",
			},
		},
	})
}
