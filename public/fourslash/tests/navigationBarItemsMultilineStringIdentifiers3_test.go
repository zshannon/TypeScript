package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestNavigationBarItemsMultilineStringIdentifiers3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `declare module 'MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters\
MoreThanOneHundredAndFiftyCharacters' { }`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineDocumentSymbol(t)
}
