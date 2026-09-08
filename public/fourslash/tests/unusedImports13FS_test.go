package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestUnusedImports13FS(t *testing.T) {
	t.Skip("Known failing fourslash test")
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @noUnusedLocals: true
// @Filename: file2.ts
[| import A, { x } from './a'; |]
console.log(A);
// @Filename: file1.ts
export default 10;
export var x = 10;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyRangeAfterCodeFix(t, `import A from './a';`, false, 0, 0)
}
