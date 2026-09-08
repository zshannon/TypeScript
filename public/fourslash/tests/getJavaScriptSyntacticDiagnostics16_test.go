package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestGetJavaScriptSyntacticDiagnostics16(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @allowJs: true
// @Filename: a.js
function F(p?) { }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineNonSuggestionDiagnostics(t)
}
