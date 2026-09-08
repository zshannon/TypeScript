package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestInlayHintsTypeMatchesName(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type Client = {};
function getClient(): Client { return {}; };
const client = getClient();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineInlayHints(t, nil /*span*/, &lsutil.UserPreferences{InlayHints: lsutil.InlayHintsPreferences{IncludeInlayVariableTypeHints: core.TSTrue, IncludeInlayVariableTypeHintsWhenTypeMatchesName: core.TSFalse}})
}
