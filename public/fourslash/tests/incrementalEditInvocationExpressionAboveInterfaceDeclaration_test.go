package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestIncrementalEditInvocationExpressionAboveInterfaceDeclaration(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
declare function alert(message?: any): void;
/*1*/
interface Foo {
    setISO8601(dString): Date;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "alert(")
	f.VerifySignatureHelp(t, fourslash.VerifySignatureHelpOptions{Text: "alert(message?: any): void"})
	f.VerifyErrorExistsAfterMarker(t, "1")
}
