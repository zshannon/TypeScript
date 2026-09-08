package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSignatureHelpTokenCrash(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
function foo(a: any, b: any) {

}

foo((/*1*/

/** This is a JSDoc comment */
foo/** More comments*/((/*2*/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifySignatureHelpWithCases(t, &fourslash.SignatureHelpCase{
		MarkerInput: "1",
		Expected:    nil,
		Context: &lsproto.SignatureHelpContext{
			IsRetrigger:      false,
			TriggerCharacter: new("("),
			TriggerKind:      lsproto.SignatureHelpTriggerKindTriggerCharacter,
		},
	})
	f.VerifySignatureHelpWithCases(t, &fourslash.SignatureHelpCase{
		MarkerInput: "2",
		Expected:    nil,
		Context: &lsproto.SignatureHelpContext{
			IsRetrigger:      false,
			TriggerCharacter: new("("),
			TriggerKind:      lsproto.SignatureHelpTriggerKindTriggerCharacter,
		},
	})
}
