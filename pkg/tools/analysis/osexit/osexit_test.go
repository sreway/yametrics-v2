package osexit_test

import (
	"testing"

	"github.com/sreway/yametrics-v2/pkg/tools/analysis/osexit"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, osexit.Analyzer, "a") // loads testdata/src/a/a.go
}
