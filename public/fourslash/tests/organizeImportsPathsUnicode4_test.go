package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestOrganizeImportsPathsUnicode4(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import * as Ab from "./Ab";
import * as _aB from "./_aB";
import * as aB from "./aB";
import * as _Ab from "./_Ab";

console.log(_aB, _Ab, aB, Ab);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import * as _Ab from "./_Ab";
import * as _aB from "./_aB";
import * as Ab from "./Ab";
import * as aB from "./aB";

console.log(_aB, _Ab, aB, Ab);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase: core.TSFalse,
			OrganizeImportsCollation:  lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsCaseFirst:  lsutil.OrganizeImportsCaseFirstUpper,
		},
	)
	f.VerifyOrganizeImports(t,
		`import * as _aB from "./_aB";
import * as _Ab from "./_Ab";
import * as aB from "./aB";
import * as Ab from "./Ab";

console.log(_aB, _Ab, aB, Ab);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase: core.TSFalse,
			OrganizeImportsCollation:  lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsCaseFirst:  lsutil.OrganizeImportsCaseFirstLower,
		},
	)
}
