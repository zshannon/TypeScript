package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestInlayHintsReparsedNodeCrash(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @allowJs: true
// @checkJs: true

// @Filename: /a.js
module.exports = function () {
  return 1;
};
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineInlayHints(t, nil /*span*/, &lsutil.UserPreferences{InlayHints: lsutil.InlayHintsPreferences{
		IncludeInlayFunctionLikeReturnTypeHints: core.TSTrue,
		IncludeInlayFunctionParameterTypeHints:  core.TSTrue,
	}})
}
