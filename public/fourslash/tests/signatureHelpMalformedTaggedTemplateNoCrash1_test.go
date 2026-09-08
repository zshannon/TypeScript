package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSignatureHelpMalformedTaggedTemplateNoCrash1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = "`${1}\n/*m1*/\n// ``\n"

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "m1")
	f.VerifyNoSignatureHelpWithContext(t, &lsproto.SignatureHelpContext{
		TriggerKind: lsproto.SignatureHelpTriggerKindInvoked,
		IsRetrigger: false,
	})
}
