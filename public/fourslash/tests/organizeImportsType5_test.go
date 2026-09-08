package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestOrganizeImportsType5(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
    d, 
    type d as D,
    type c,
    c as C,
    b,
    b as B,
    type A,
    a
} from './foo';
console.log(A, a, B, b, c, C, d, D);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import {
    type A,
    a,
    b,
    b as B,
    type c,
    c as C,
    d,
    type d as D
} from './foo';
console.log(A, a, B, b, c, C, d, D);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase: core.TSUnknown,
			OrganizeImportsTypeOrder:  lsutil.OrganizeImportsTypeOrderInline,
		},
	)
	f.VerifyOrganizeImports(t,
		`import {
    type A,
    a,
    b,
    b as B,
    type c,
    c as C,
    d,
    type d as D
} from './foo';
console.log(A, a, B, b, c, C, d, D);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase: core.TSUnknown,
			OrganizeImportsTypeOrder:  lsutil.OrganizeImportsTypeOrderInline,
		},
	)
}
