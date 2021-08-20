package capacity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-cache/capacity"
)

var _ = Describe("MaxLen", func() {
	var (
		mockLen uint
		subject capacity.Tracker
	)
	When("no capacity", func() {
		BeforeEach(func() {
			subject = capacity.NewMaxLen(0, func() uint {
				return mockLen
			})
		})
		It("can't fit anything", func() {
			Expect(subject.IsLargerThanCapacity(1)).Should(BeTrue())
		})
		It("can't fit", func() {
			Expect(subject.HasCapacity(1)).Should(BeFalse())
		})
	})

	When("non-zero capacity", func() {
		max := uint(5)
		BeforeEach(func() {
			subject = capacity.NewMaxLen(max, func() uint {
				return mockLen
			})
		})
		When("empty", func() {
			It("can fit items smaller than capacity", func() {
				Expect(subject.IsLargerThanCapacity(1)).Should(BeFalse())
			})
			It("can fit item exactly capacity", func() {
				Expect(subject.IsLargerThanCapacity(max)).Should(BeFalse())
			})
			It("can't fit items larger than capacity", func() {
				Expect(subject.IsLargerThanCapacity(max + 1)).Should(BeTrue())
			})
		})
		When("not empty", func() {
			BeforeEach(func() {
				mockLen = 3
			})
			It("can fit items small enough", func() {
				Expect(subject.HasCapacity(1)).Should(BeTrue())
			})
			It("can't fit items that are too large", func() {
				Expect(subject.HasCapacity(4)).Should(BeFalse())
			})
		})
		When("full", func() {
			BeforeEach(func() {
				mockLen = max
			})
			It("can't fit anything'", func() {
				Expect(subject.HasCapacity(1)).Should(BeFalse())
			})
		})
	})
})
