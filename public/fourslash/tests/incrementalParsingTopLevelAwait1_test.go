package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestIncrementalParsingTopLevelAwait1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @target: esnext
// @module: esnext
// @Filename: ./foo.ts
await(1);
/*1*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyNumberOfErrorsInCurrentFile(t, 1)
	f.GoToMarker(t, "1")
	f.Insert(t, "export {};")
	f.VerifyNumberOfErrorsInCurrentFile(t, 0)
	f.ReplaceLine(t, 1, "")
	f.VerifyNumberOfErrorsInCurrentFile(t, 1)
}
