package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocPropertyDescription1(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface StringExample {
    /** Something generic */
    [p: string]: any; 
    /** Something specific */
    property: number;
}
function stringExample(e: StringExample) {
    console.log(e./*property*/property);
    console.log(e./*string*/anything); 
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "property", "(property) StringExample.property: number", "Something specific")
	f.VerifyQuickInfoAt(t, "string", "(index) StringExample[string]: any", "Something generic")
}
