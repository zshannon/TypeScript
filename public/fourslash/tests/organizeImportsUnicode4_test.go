package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestOrganizeImportsUnicode4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import {
    Ab,
    _aB,
    aB,
    _Ab,
} from './foo';

console.log(_aB, _Ab, aB, Ab);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import {
    _Ab,
    _aB,
    Ab,
    aB,
} from './foo';

console.log(_aB, _Ab, aB, Ab);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase: core.TSFalse,
			OrganizeImportsCollation:  lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsCaseFirst:  lsutil.OrganizeImportsCaseFirstUpper,
		},
	)
	f.VerifyOrganizeImports(t,
		`import {
    _aB,
    _Ab,
    aB,
    Ab,
} from './foo';

console.log(_aB, _Ab, aB, Ab);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase: core.TSFalse,
			OrganizeImportsCollation:  lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsCaseFirst:  lsutil.OrganizeImportsCaseFirstLower,
		},
	)
}
