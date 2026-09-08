package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestJsDocSee4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `class [|/*def1*/A|] {
    foo () { }
}
declare const [|/*def2*/a|]: A;
/**
 * @see {/*use1*/[|A|]#foo}
 */
const t1 = 1
/**
 * @see {/*use2*/[|a|].foo()}
 */
const t2 = 1
/**
 * @see {@link /*use3*/[|a|].foo()}
 */
const t3 = 1`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToDefinition(t, false, "use1", "use2", "use3")
}
