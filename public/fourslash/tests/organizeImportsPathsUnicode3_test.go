package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestOrganizeImportsPathsUnicode3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import * as B from "./B";
import * as À from "./À";
import * as A from "./A";

console.log(A, À, B);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import * as A from "./A";
import * as À from "./À";
import * as B from "./B";

console.log(A, À, B);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase:      core.TSFalse,
			OrganizeImportsCollation:       lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsAccentCollation: core.TSFalse,
		},
	)
	f.VerifyOrganizeImports(t,
		`import * as A from "./A";
import * as À from "./À";
import * as B from "./B";

console.log(A, À, B);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		&lsutil.UserPreferences{
			OrganizeImportsIgnoreCase:      core.TSFalse,
			OrganizeImportsCollation:       lsutil.OrganizeImportsCollationUnicode,
			OrganizeImportsAccentCollation: core.TSTrue,
		},
	)
}
