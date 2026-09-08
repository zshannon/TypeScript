package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestDocCommentTemplateInMultiLineComment(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: justAComment.ts
/* /*0*/ */`
	capabilities := fourslash.GetDefaultCapabilities()
	capabilities.TextDocument.Completion.CompletionItem.SnippetSupport = new(false)
	f, done := fourslash.NewFourslash(t, capabilities, content)
	defer done()
	f.VerifyNoJSDocCompletion(t, "0")
}
