package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestOrganizeImportsUnicode2(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
    a2,
    a100,
    a1,
} from './foo';

console.log(a1, a2, a100);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import {
    a1,
    a100,
    a2,
} from './foo';

console.log(a1, a2, a100);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase:       core.TSFalse,
			OrganizeImportsCollation:        lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsNumericCollation: core.TSFalse,
		},
	)
	f.VerifyOrganizeImports(t,
		`import {
    a1,
    a2,
    a100,
} from './foo';

console.log(a1, a2, a100);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase:       core.TSFalse,
			OrganizeImportsCollation:        lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsNumericCollation: core.TSTrue,
		},
	)
}
