package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestInlayHintsQuotePreference1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const a1: '"' = '"';
const b1: '\\' = '\\';
export function fn(a = a1, b = b1) {}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineInlayHints(t, nil /*span*/, &lsutil.UserPreferences{QuotePreference: lsutil.QuotePreference("double"), InlayHints: lsutil.InlayHintsPreferences{IncludeInlayFunctionParameterTypeHints: core.TSTrue}})
}
