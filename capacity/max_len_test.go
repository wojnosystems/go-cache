package capacity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-cache/capacity"
)

var _ = Describe("MaxLen", func() {
	var (
		subject capacity.TrackMutator
	)
	When("no capacity", func() {
		BeforeEach(func() {
			subject = capacity.NewMaxLen(0)
		})
		It("can't fit anything", func() {
			Expect(subject.IsLargerThanCapacity(1)).Should(BeTrue())
		})
		It("can't fit", func() {
			Expect(subject.Add(1)).Should(BeFalse())
		})
	})

	When("non-zero capacity", func() {
		max := uint(5)
		BeforeEach(func() {
			subject = capacity.NewMaxLen(max)
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
			It("adds item", func() {
				Expect(subject.Add(max)).Should(BeTrue())
			})
			It("does not underflow", func() {
				subject.Remove(max + 100)
				Expect(subject.Add(max)).Should(BeTrue())
				Expect(subject.Add(max)).Should(BeFalse())
			})
		})
		When("not empty", func() {
			var (
				startLen uint
			)
			BeforeEach(func() {
				startLen = 3
				subject.Add(startLen)
			})
			It("can fit items small enough", func() {
				Expect(subject.Add(1)).Should(BeTrue())
			})
			It("can't fit items that are too large", func() {
				Expect(subject.Add(4)).Should(BeFalse())
			})
			It("removes items", func() {
				subject.Remove(1)
				Expect(subject.Add(max - (startLen - 1))).Should(BeTrue())
			})
		})
		When("full", func() {
			BeforeEach(func() {
				subject.Add(max)
			})
			It("can't fit anything'", func() {
				Expect(subject.Add(1)).Should(BeFalse())
			})
			It("won't add the value'", func() {
				Expect(subject.Add(1)).Should(BeFalse())
			})
		})
	})
})
