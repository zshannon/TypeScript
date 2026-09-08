package fourslash_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/fourslash"
	"github.com/zshannon/TypeScript/public/v7/testutil"
)

func TestCodefixClassImplementInterface_omit(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `interface One {
    a: number;
    b: string;
}

interface Two extends Omit<One, "a"> {
    c: boolean;
}

class TwoStore implements Two {[| |]}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyCodeFix(t, fourslash.VerifyCodeFixOptions{
		Description: "Implement interface 'Two'",
		NewFileContent: `interface One {
    a: number;
    b: string;
}

interface Two extends Omit<One, "a"> {
    c: boolean;
}

class TwoStore implements Two {
    c: boolean;
    b: string;
}`,
		Index: 0,
	})
}
