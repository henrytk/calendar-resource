package calendar_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCalendar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Calendar Suite")
}
