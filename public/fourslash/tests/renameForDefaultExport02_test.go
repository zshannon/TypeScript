package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/fourslash"
	. "github.com/zshannon/TypeScript/public/fourslash/tests/util"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestRenameForDefaultExport02(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `[|export default function /*1*/[|{| "contextRangeIndex": 0 |}DefaultExportedFunction|]() {
    return /*2*/[|DefaultExportedFunction|]
}|]
/**
 *  Commenting [|{| "inComment": true |}DefaultExportedFunction|]
 */

var x: typeof /*3*/[|DefaultExportedFunction|];

var y = /*4*/[|DefaultExportedFunction|]();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineRename(t, nil /*preferences*/, ToAny(core.Filter(f.GetRangesByText().Get("DefaultExportedFunction"), func(r *fourslash.RangeMarker) bool { return r.Marker == nil || r.Marker.Data["inComment"] == nil }))...)
}
