package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestSquiggleIllegalClassExtension(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class Foo extends /*1*/Bar/*2*/ { }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyErrorExistsBetweenMarkers(t, "1", "2")
	f.VerifyNumberOfErrorsInCurrentFile(t, 1)
}
