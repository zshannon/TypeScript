package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocSignature_43394(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `/**
 * @typedef {Object} Foo
 * @property {number} ...
 * /**/@typedef {number} Bar
 */`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineSignatureHelp(t)
}
