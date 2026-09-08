package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestNavigateToSymbolIterator(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @lib: es5
class C {
    [|[Symbol.iterator]|]() {}
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyWorkspaceSymbol(t, []*fourslash.VerifyWorkspaceSymbolCase{
		{
			Pattern:     "iterator",
			Preferences: nil,
			Exact: new([]*lsproto.SymbolInformation{
				{
					Name:          "iterator",
					Kind:          lsproto.SymbolKindMethod,
					Location:      f.Ranges()[0].LSLocation(),
					ContainerName: new("C"),
				},
			}),
		},
	})
}
