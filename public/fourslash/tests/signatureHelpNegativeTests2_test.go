package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSignatureHelpNegativeTests2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class clsOverload { constructor(); constructor(test: string); constructor(test?: string) { } }
var x = new clsOverload/*beforeOpenParen*/()/*afterCloseParen*/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyNoSignatureHelpForMarkers(t, "beforeOpenParen", "afterCloseParen")
}
