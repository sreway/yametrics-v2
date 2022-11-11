package osexit_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sreway/yametrics-v2/pkg/tools/analysis/osexit"
)

func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, osexit.Analyzer, "a", "b") // loads testdata/src/
}
