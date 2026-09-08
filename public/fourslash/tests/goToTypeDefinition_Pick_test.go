package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestGoToTypeDefinition_Pick(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
type User = { id: number; name: string; };
declare const user: Pick<User, "name">
/*reference*/user

type PickedUser = Pick<User, "name">
declare const user2: PickedUser
/*reference2*/user2`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToTypeDefinition(t, "reference", "reference2")
}
