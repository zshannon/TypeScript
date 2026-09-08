package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestQuickInfoSignatureRestParameterFromUnion2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strict: false
declare const rest:
  | ((a?: { a: true }, ...rest: string[]) => unknown)
  | ((b?: { b: true }) => unknown);

/**/rest({ a: true, b: true }, "foo", "bar");`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "", "const rest: (arg0?: {\n    a: true;\n} & {\n    b: true;\n}, ...rest: string[]) => unknown", "")
}
