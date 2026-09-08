package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestDocCommentTemplateWithExistingJSDoc(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/** /**/ */

/**
 * @param {string} a
 * @param {string} b
 */
function foo(a, b) {
    return a + b;
}`
	capabilities := fourslash.GetDefaultCapabilities()
	capabilities.TextDocument.Completion.CompletionItem.SnippetSupport = new(false)
	f, done := fourslash.NewFourslash(t, capabilities, content)
	defer done()
	f.VerifyNoJSDocCompletion(t, "")
}
