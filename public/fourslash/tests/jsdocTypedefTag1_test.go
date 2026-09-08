package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsdocTypedefTag1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es2015
// @allowNonTsExtensions: true
// @Filename: jsdocCompletion_typedef.js
/**
 * @typedef {Object} MyType
 * @property {string} yes
 */
function foo() { }
/**
 * @param {MyType} my
 */
function a(my) {
    my.yes./*1*/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.MarkTestAsStradaServer()
	f.VerifyCompletions(t, "1", &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Includes: []fourslash.CompletionsExpectedItem{
				"charAt",
			},
		},
	})
}
