package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/fourslash"
	"github.com/zshannon/TypeScript/public/lsp/lsproto"
	"github.com/zshannon/TypeScript/public/testutil"
)

func TestGetOutliningSpansForImports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `[|import * as ns from "mod";

import d from "mod";
import { a, b, c } from "mod";

import r = require("mod");|]

// statement
var x = 0;

// another set of imports
[|import * as ns from "mod";
import d from "mod";
import { a, b, c } from "mod";
import r = require("mod");|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOutliningSpans(t, lsproto.FoldingRangeKindImports)
}
