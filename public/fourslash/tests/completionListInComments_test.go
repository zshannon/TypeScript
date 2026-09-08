package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestCompletionListInComments(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `var foo = '';
( // f/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCompletions(t, "", nil)
}
