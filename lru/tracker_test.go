package lru_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-cache/lru"
)

var _ = Describe("Tracker", func() {
	var (
		subject lru.Tracker
	)
	BeforeEach(func() {
		subject = lru.NewTracker()
	})

	When("empty", func() {
		It("has no lru", func() {
			_, ok := subject.LRU()
			Expect(ok).Should(BeFalse())
		})
		It("is empty", func() {
			Expect(subject.Len()).Should(BeZero())
		})
		It("removes nothing", func() {
			subject.Remove(1)
		})
	})

	When("one item", func() {
		itemKey := 1
		BeforeEach(func() {
			subject.Touch(itemKey)
		})
		It("updates the length", func() {
			Expect(subject.Len()).Should(Equal(1))
		})
		It("is the LRU", func() {
			actual, ok := subject.LRU()
			Expect(ok).Should(BeTrue())
			Expect(actual).Should(Equal(itemKey))
		})
		When("removed", func() {
			It("is empty", func() {
				subject.Remove(itemKey)
				Expect(subject.Len()).Should(BeZero())
			})
		})
	})

	When("multiple items", func() {
		BeforeEach(func() {
			subject.Touch(1)
			subject.Touch(2)
			subject.Touch(3)
			subject.Touch(4)
		})
		It("updates the length", func() {
			Expect(subject.Len()).Should(Equal(4))
		})
		It("tracks the LRU", func() {
			actual, _ := subject.LRU()
			Expect(actual).Should(Equal(1))
		})
		When("item is touched", func() {
			It("updates the LRU", func() {
				subject.Touch(1)
				actual, _ := subject.LRU()
				Expect(actual).Should(Equal(2))
			})
			It("updates the LRU each time", func() {
				subject.Touch(1)
				subject.Touch(2)
				subject.Touch(3)
				actual, _ := subject.LRU()
				Expect(actual).Should(Equal(4))
			})
		})
		When("item is removed", func() {
			BeforeEach(func() {
				subject.Remove(1)
			})
			It("updates the LRU", func() {
				actual, _ := subject.LRU()
				Expect(actual).Should(Equal(2))
			})
			It("reduces the length", func() {
				Expect(subject.Len()).Should(Equal(3))
			})
		})
		When("existing item is touched", func() {
			It("does not change the length", func() {
				before := subject.Len()
				subject.Touch(1)
				subject.Touch(2)
				subject.Touch(3)
				subject.Touch(4)
				Expect(subject.Len()).Should(Equal(before))
			})
		})
	})
})
