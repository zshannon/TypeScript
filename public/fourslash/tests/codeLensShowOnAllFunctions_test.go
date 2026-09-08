package fourslash_test

import (
	"fmt"
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/ls/lsutil"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodeLensReferencesShowOnAllFunctions(t *testing.T) {
	t.Parallel()
	containingTestName := t.Name()
	for _, value := range []core.Tristate{core.TSTrue, core.TSFalse} {
		t.Run(fmt.Sprintf("%s=%v", containingTestName, value.IsTrue()), func(t *testing.T) {
			t.Parallel()
			defer testutil.RecoverAndFail(t, "Panic on fourslash test")

			const content = `
export function f1(): void {}

function f2(): void {}

export const f3 = () => {};

const f4 = () => {};

const f5 = function() {};
`
			f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
			defer done()
			f.VerifyBaselineCodeLens(t, &lsutil.UserPreferences{
				CodeLens: lsutil.CodeLensUserPreferences{
					ReferencesCodeLensEnabled:            core.TSTrue,
					ReferencesCodeLensShowOnAllFunctions: value,
				},
			})
		})
	}
}
