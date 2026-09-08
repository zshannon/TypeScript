package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestFixExactOptionalUnassignableProperties3(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @strictNullChecks: true
// @exactOptionalPropertyTypes: true
// @Filename: fixExactOptionalUnassignableProperties2.ts
import { INodeModules } from 'foo'
interface J {
    a?: number | undefined
}
declare var inm: INodeModules
declare var j: J
inm/**/ = j
console.log(inm)
// @Filename: node_modules/@types/foo/index.d.ts
export interface INodeModules {
    a?: number
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCodeFixNotAvailable(t)
}
