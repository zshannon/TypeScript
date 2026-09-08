package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestFormatSelectionPreserveTrailingWhitespace(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
/*begin*/;    
    
/*end*/    
    
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	opts154 := f.GetOptions()
	opts154.FormatCodeSettings.TrimTrailingWhitespace = core.TSFalse
	f.Configure(t, opts154)
	f.FormatSelection(t, "begin", "end")
	f.VerifyCurrentFileContent(t, `
;    
    
    
    
`)
}
