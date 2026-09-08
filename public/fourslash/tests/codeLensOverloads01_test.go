package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodeLensOverloads01(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `
export function foo(x: number): number;
export function foo(x: string): string;
export function foo(x: string | number): string | number {
	return x;
}

foo(1);

foo("hello");

// This one isn't expected to match any overload,
// but is really just here to test how it affects how code lens.
foo(Math.random() ? 1 : "hello");
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineCodeLens(t, &lsutil.UserPreferences{
		CodeLens: lsutil.CodeLensUserPreferences{
			ReferencesCodeLensEnabled:            core.TSTrue,
			ReferencesCodeLensShowOnAllFunctions: core.TSTrue,

			ImplementationsCodeLensEnabled:                core.TSTrue,
			ImplementationsCodeLensShowOnInterfaceMethods: core.TSTrue,
			ImplementationsCodeLensShowOnAllClassMethods:  core.TSTrue,
		},
	})
}
