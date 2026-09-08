package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGoToTypeDefinitionUnionType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class /*definition0*/C {
    p;
}

interface /*definition1*/I {
    x;
}

namespace M {
    export interface /*definition2*/I {
        y;
    }
}

var x: C | I | M.I;

/*reference*/x;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToTypeDefinition(t, "reference")
}
