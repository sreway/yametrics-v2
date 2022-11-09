package osexit_test

import (
	"github.com/sreway/yametrics-v2/pkg/tools/analysis/osexit"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestFromFileSystem(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, osexit.Analyzer, "a") // loads testdata/src/a/a.go
}
