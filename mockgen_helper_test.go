package cache_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
)

// ginkgoTestReporter tricks gomock to use Ginkgo reporting instead of Go's native reporting
// https://github.com/onsi/ginkgo/issues/9
type ginkgoTestReporter struct{}

func (g ginkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g ginkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}
