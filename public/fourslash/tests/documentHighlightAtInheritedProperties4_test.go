package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestDocumentHighlightAtInheritedProperties4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: file1.ts
class class1 extends class1 {
   [|doStuff|]() { }
   [|propName|]: string;
}

var c: class1;
c.[|doStuff|]();
c.[|propName|];`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineDocumentHighlights(t, nil /*preferences*/, ToAny(f.Ranges())...)
}
