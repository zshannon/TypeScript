package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestDocumentHighlightVarianceModifiers(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type TFoo<Value> = { value: Value };
type TBar<[|in|] [|out|] Value> = TFoo<Value>;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineDocumentHighlights(t, nil /*preferences*/, ToAny(f.Ranges())...)
}
