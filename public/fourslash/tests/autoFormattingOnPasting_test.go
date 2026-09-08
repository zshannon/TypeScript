package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestAutoFormattingOnPasting(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `namespace TestModule {
/**/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Paste(t, " class TestClass{\nprivate   foo;\npublic testMethod( )\n{}\n}")
	f.VerifyCurrentFileContent(t, `namespace TestModule {
    class TestClass {
        private foo;
        public testMethod() { }
    }
}`)
}
