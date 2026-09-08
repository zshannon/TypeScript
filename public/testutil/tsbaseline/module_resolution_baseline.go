package tsbaseline

import (
	"testing"

	"github.com/zshannon/TypeScript/public/v7/testutil/baseline"
)

func DoModuleResolutionBaseline(t *testing.T, baselinePath string, trace string, opts baseline.Options) {
	baselinePath = tsExtension.ReplaceAllString(baselinePath, ".trace.json")
	var errorBaseline string
	if trace != "" {
		errorBaseline = trace
	} else {
		errorBaseline = baseline.NoContent
	}
	baseline.Run(t, baselinePath, errorBaseline, opts)
}
