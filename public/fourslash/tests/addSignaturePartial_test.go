package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestAddSignaturePartial(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = ``
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.Insert(t, "interface Number { toFixed")
	f.Insert(t, "(")
}
