package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestFormatInsertSpaceAfterCloseBraceBeforeCloseBracket(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `[{}]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	opts122 := f.GetOptions()
	opts122.FormatCodeSettings.InsertSpaceAfterOpeningAndBeforeClosingNonemptyBrackets = core.TSTrue
	f.Configure(t, opts122)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `[ {} ]`)
}
