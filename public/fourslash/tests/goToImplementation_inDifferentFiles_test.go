package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGoToImplementation_inDifferentFiles(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
// @Filename: /home/src/workspaces/project/bar.ts
import {Foo} from './foo'

class [|A|] implements Foo {
    func() {}
}

class [|B|] implements Foo {
    func() {}
}
// @Filename: /home/src/workspaces/project/foo.ts
export interface /**/Foo {
    func();
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.MarkTestAsStradaServer()
	f.VerifyBaselineGoToImplementation(t, "")
}
