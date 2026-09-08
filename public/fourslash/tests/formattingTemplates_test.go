package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFormattingTemplates(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `String.call ` + "`" + `${123}` + "`" + `/*1*/
String.call ` + "`" + `${123} ${456}` + "`" + `/*2*/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, ";")
	f.VerifyCurrentLineContent(t, "String.call`${123}`;")
	f.GoToMarker(t, "2")
	f.Insert(t, ";")
	f.VerifyCurrentLineContent(t, "String.call`${123} ${456}`;")
}
