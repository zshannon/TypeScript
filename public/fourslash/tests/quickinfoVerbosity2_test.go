package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestQuickinfoVerbosity2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type Str = string | {};
type FooType = Str | number;
type Sym = symbol | (() => void);
type BarType = Sym | boolean;
type BothType = FooType | BarType;
const both/*b*/: BothType = 1;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineHoverWithVerbosity(t, map[string][]int{"b": {0, 1, 2, 3}})
}
