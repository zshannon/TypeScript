package encoder_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/core"
	"github.com/zshannon/TypeScript/public/testutil/baseline"
)

func TestMain(m *testing.M) {
	core.ApplyDebugStackLimit()
	defer baseline.Track()()
	m.Run()
}
