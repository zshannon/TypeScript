package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestVerifySingleFileEmitOutput1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: verifySingleFileEmitOutput1_file0.ts
export class A {
}
export class Z {
}
// @Filename: verifySingleFileEmitOutput1_file1.ts
import f = require("./verifySingleFileEmitOutput1_file0");
var /**/b = new f.A();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "", "var b: f.A", "")
}
