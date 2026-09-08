package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodeFixInferFromCallInAssignment(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noImplicitAny: true
function inferAny( [| app |] ) {
    const result = app.use('hi')
    return result
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyRangeAfterCodeFix(t, `app: { use: (arg0: string) => any }`, false, 0, 0)
}
