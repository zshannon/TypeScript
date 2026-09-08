package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/ls/lsutil"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestImportNameCodeFix_typesVersions(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: commonjs
// @checkJs: true
// @Filename: /node_modules/unified/package.json
{
  "name": "unified",
  "types": "types/ts3.444/index.d.ts",
  "typesVersions": {
    ">=4.0": {
      "types/ts3.444/*": [
        "types/ts4.0/*"
      ]
    }
  }
}
// @Filename: /node_modules/unified/types/ts3.444/index.d.ts
export declare const x: number;
// @Filename: /node_modules/unified/types/ts4.0/index.d.ts
export declare const x: number;
// @Filename: /foo.js
import {} from "unified";
// @Filename: /index.js
x/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "", []string{"unified", "unified/types/ts3.444/index.js"}, &lsutil.UserPreferences{ImportModuleSpecifierEnding: "js", AutoImportEntrypointDirectorySearch: core.TSTrue})
}
