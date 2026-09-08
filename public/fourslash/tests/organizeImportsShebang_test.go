package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestOrganizeImports_Shebang_PreserveAndSort(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `#!/usr/bin/env node
import Foo from "foo";
import Bar from "bar";

import Foobar from "foobar";

console.log(Foo, Bar, Foobar);`

	f, done := fourslash.NewFourslash(t, nil /* capabilities */, content)
	defer done()

	f.VerifyOrganizeImports(
		t,
		`#!/usr/bin/env node
import Bar from "bar";
import Foo from "foo";

import Foobar from "foobar";

console.log(Foo, Bar, Foobar);`,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		nil,
	)
}
