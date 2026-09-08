package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCallHierarchyConstNamedClassExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo() {
    new Bar();
}

const /**/Bar = class {
    constructor() {
        baz();
    }
}

function baz() {
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.VerifyBaselineCallHierarchy(t)
}
