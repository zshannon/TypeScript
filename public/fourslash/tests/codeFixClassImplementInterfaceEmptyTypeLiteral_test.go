package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestCodeFixClassImplementInterfaceEmptyTypeLiteral(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
interface I {
    x: {};
}

class C implements I {[|
   |]constructor() { }
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCodeFix(t, fourslash.VerifyCodeFixOptions{
		Description: "Implement interface 'I'",
		NewFileContent: `
interface I {
    x: {};
}

class C implements I {
   constructor() { }
    x: {};
}`,
		Index: 0,
	})
}
