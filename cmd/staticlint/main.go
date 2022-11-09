package main

import (
	"strings"

	"github.com/sreway/yametrics-v2/pkg/tools/analysis/osexit"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	analyzers := []*analysis.Analyzer{}
	analyzers = append(analyzers, asmdecl.Analyzer)
	analyzers = append(analyzers, assign.Analyzer)
	analyzers = append(analyzers, atomic.Analyzer)
	analyzers = append(analyzers, atomicalign.Analyzer)
	analyzers = append(analyzers, bools.Analyzer)
	analyzers = append(analyzers, buildssa.Analyzer)
	analyzers = append(analyzers, buildtag.Analyzer)
	analyzers = append(analyzers, cgocall.Analyzer)
	analyzers = append(analyzers, composite.Analyzer)
	analyzers = append(analyzers, copylock.Analyzer)
	analyzers = append(analyzers, ctrlflow.Analyzer)
	analyzers = append(analyzers, deepequalerrors.Analyzer)
	analyzers = append(analyzers, errorsas.Analyzer)
	analyzers = append(analyzers, fieldalignment.Analyzer)
	analyzers = append(analyzers, findcall.Analyzer)
	analyzers = append(analyzers, framepointer.Analyzer)
	analyzers = append(analyzers, httpresponse.Analyzer)
	analyzers = append(analyzers, ifaceassert.Analyzer)
	analyzers = append(analyzers, inspect.Analyzer)
	analyzers = append(analyzers, loopclosure.Analyzer)
	analyzers = append(analyzers, lostcancel.Analyzer)
	analyzers = append(analyzers, nilfunc.Analyzer)
	analyzers = append(analyzers, nilness.Analyzer)
	analyzers = append(analyzers, pkgfact.Analyzer)
	analyzers = append(analyzers, printf.Analyzer)
	analyzers = append(analyzers, reflectvaluecompare.Analyzer)
	analyzers = append(analyzers, shadow.Analyzer)
	analyzers = append(analyzers, shift.Analyzer)
	analyzers = append(analyzers, sigchanyzer.Analyzer)
	analyzers = append(analyzers, sortslice.Analyzer)
	analyzers = append(analyzers, stdmethods.Analyzer)
	analyzers = append(analyzers, stringintconv.Analyzer)
	analyzers = append(analyzers, structtag.Analyzer)
	analyzers = append(analyzers, testinggoroutine.Analyzer)
	analyzers = append(analyzers, tests.Analyzer)
	analyzers = append(analyzers, unmarshal.Analyzer)
	analyzers = append(analyzers, unreachable.Analyzer)
	analyzers = append(analyzers, unsafeptr.Analyzer)
	analyzers = append(analyzers, unusedresult.Analyzer)
	analyzers = append(analyzers, unusedwrite.Analyzer)

	for _, i := range staticcheck.Analyzers {
		if strings.HasPrefix(i.Analyzer.Name, "SA") {
			analyzers = append(analyzers, i.Analyzer)
		}
	}

	for _, i := range simple.Analyzers {
		analyzers = append(analyzers, i.Analyzer)
	}

	analyzers = append(analyzers, osexit.Analyzer)

	multichecker.Main(analyzers...)
}
