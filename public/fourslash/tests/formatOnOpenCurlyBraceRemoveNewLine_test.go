package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestFormatOnOpenCurlyBraceRemoveNewLine(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `if(true)
/**/ }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	opts124 := f.GetOptions()
	opts124.FormatCodeSettings.PlaceOpenBraceOnNewLineForControlBlocks = core.TSFalse
	f.Configure(t, opts124)
	f.GoToMarker(t, "")
	f.Insert(t, "{")
	f.VerifyCurrentFileContent(t, `if (true) { }`)
}
