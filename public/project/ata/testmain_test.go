package ata_test

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/core"
	"github.com/zshannon/TypeScript/public/v7/testutil/baseline"
)

func TestMain(m *testing.M) {
	core.ApplyDebugStackLimit()
	defer baseline.Track()()
	m.Run()
}
