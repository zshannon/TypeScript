package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestRenameStringPropertyNames2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `type Props = {
  foo: boolean;
}

let { foo }: Props = null as any;
foo;

let asd: Props = { "foo"/**/: true }; // rename foo here`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineRename(t, nil /*preferences*/, "")
}
