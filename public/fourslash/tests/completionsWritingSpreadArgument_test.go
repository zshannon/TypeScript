package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	. "github.com/zshannon/TypeScript/public/v7/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCompletionsWritingSpreadArgument(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5

const [] = [Math.min(./*marker*/)]
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyCompletions(t, nil, nil)
	f.Insert(t, ".")
	f.VerifyCompletions(t, nil, nil)
	f.Insert(t, ".")
	f.VerifyCompletions(t, nil, &fourslash.CompletionsExpectedList{
		IsIncomplete: false,
		ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
			CommitCharacters: &DefaultCommitCharacters,
			EditRange:        Ignored,
		},
		Items: &fourslash.CompletionsExpectedItems{
			Exact: CompletionGlobals,
		},
	})
}
