package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodeFixClassImplementInterfaceOptionalProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strict: false
interface IPerson {
    name: string;
    birthday?: string;
}
class Person implements IPerson {}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCodeFix(t, fourslash.VerifyCodeFixOptions{
		Description: "Implement interface 'IPerson'",
		NewFileContent: `interface IPerson {
    name: string;
    birthday?: string;
}
class Person implements IPerson {
    name: string;
    birthday?: string;
}`,
		Index: 0,
	})
}
