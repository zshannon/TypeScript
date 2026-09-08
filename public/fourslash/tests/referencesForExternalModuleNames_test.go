package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestReferencesForExternalModuleNames(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: referencesForGlobals_1.ts
/*1*/declare module "/*2*/foo" {
    var f: number;
}
// @Filename: referencesForGlobals_2.ts
/*3*/import f = require("/*4*/foo");`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineFindAllReferences(t, "1", "2", "3", "4")
}
