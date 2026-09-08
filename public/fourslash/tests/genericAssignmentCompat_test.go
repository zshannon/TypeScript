package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGenericAssignmentCompat(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface Int<T> {

    val<U>(f: (t: T) => U): Int<U>;

}

declare var v1: Int<string>;

var /*1*/v2/*2*/: Int<number> = v1;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyErrorExistsBetweenMarkers(t, "1", "2")
	f.VerifyNumberOfErrorsInCurrentFile(t, 1)
}
