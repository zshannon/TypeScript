package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodeFixClassImplementInterfaceWithAmbientSignatures2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare class A {
    method(): void;
}
class B implements A {}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCodeFix(t, fourslash.VerifyCodeFixOptions{
		Description: "Implement interface 'A'",
		NewFileContent: `declare class A {
    method(): void;
}
class B implements A {
    method(): void {
        throw new Error("Method not implemented.");
    }
}`,
		Index: 0,
	})
}
