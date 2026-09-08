package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodeFixConvertToTypeOnlyImport1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @verbatimModuleSyntax: true
// @Filename: exports.ts
export default class A {}
export class B {}
export class C {}
// @Filename: imports.ts
import {
    B,
    C,
} from './exports';

declare const b: B;
declare const c: C;
console.log(b, c);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToFile(t, "imports.ts")
	f.VerifyCodeFixNotAvailable(t)
}
