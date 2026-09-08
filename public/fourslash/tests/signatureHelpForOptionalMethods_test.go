package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSignatureHelpForOptionalMethods(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strict: true
interface Obj {
    optionalMethod?: (current: any) => any;
};

const o: Obj = {
  optionalMethod(/*1*/) {
    return {};
  }
};`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.VerifySignatureHelp(t, fourslash.VerifySignatureHelpOptions{Text: "optionalMethod(current: any): any", ParameterName: "current", ParameterSpan: "current: any"})
}
