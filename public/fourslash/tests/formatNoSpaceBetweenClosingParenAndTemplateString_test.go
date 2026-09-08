package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormatNoSpaceBetweenClosingParenAndTemplateString(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `foo() ` + "`" + `abc` + "`" + `;
bar()` + "`" + `def` + "`" + `;
baz()` + "`" + `a${x}b` + "`" + `;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, "foo()`abc`;\nbar()`def`;\nbaz()`a${x}b`;")
}
