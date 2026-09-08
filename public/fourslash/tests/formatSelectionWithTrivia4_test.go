package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormatSelectionWithTrivia4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if (true) {
/*begin*/// test comment
/*end*/console.log();
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatSelection(t, "begin", "end")
	f.VerifyCurrentFileContent(t, `if (true) {
    // test comment
console.log();
}`)
}
