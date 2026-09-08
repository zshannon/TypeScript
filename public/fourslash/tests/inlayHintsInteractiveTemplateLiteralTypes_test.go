package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestInlayHintsInteractiveTemplateLiteralTypes(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare function getTemplateLiteral1(): ` + "`" + `${string},${string}` + "`" + `;
const lit1 = getTemplateLiteral1();
declare function getTemplateLiteral2(): ` + "`" + `\${${string},${string}` + "`" + `;
const lit2 = getTemplateLiteral2();
declare function getTemplateLiteral3(): ` + "`" + `start${string}\${,$${string}end` + "`" + `;
const lit3 = getTemplateLiteral3();
declare function getTemplateLiteral4(): ` + "`" + `${string}\` + "`" + `,${string}` + "`" + `;
const lit4 = getTemplateLiteral4();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineInlayHints(t, nil /*span*/, &lsutil.UserPreferences{InlayHints: lsutil.InlayHintsPreferences{IncludeInlayVariableTypeHints: core.TSTrue}})
}
