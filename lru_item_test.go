package cache_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-cache"
)

var _ = Describe("LRUItem", func() {
	var (
		ctrl    *gomock.Controller
		source  *MockstatefulValueMapper
		subject cache.GetInvalidater
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(ginkgoTestReporter{})
		source = NewMockstatefulValueMapper(ctrl)
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	When("not at capacity", func() {
		BeforeEach(func() {
			subject = cache.NewLRUItem(10, valueMapperWrap(source))
		})

		When("get", func() {
			BeforeEach(func() {
				source.EXPECT().Get(gomock.Any(), "1").Times(1).Return("1", nil)
			})
			It("is cached", func() {
				Expect(subject.Get(ignoreCtx, "1")).Should(Equal("1"))
				Expect(subject.Get(ignoreCtx, "1")).Should(Equal("1"))
			})
		})

		When("invalidated", func() {
			When("with elements", func() {
				BeforeEach(func() {
					source.EXPECT().Get(ignoreCtx, "1").Times(2).Return("1", nil)
					_, _ = subject.Get(ignoreCtx, "1")
				})

				It("fetches a new item", func() {
					subject.Invalidate("1")
					_, _ = subject.Get(ignoreCtx, "1")
				})
			})
			When("without elements", func() {
				BeforeEach(func() {
					source.EXPECT().Get(ignoreCtx, "1").Times(1).Return("1", nil)
				})

				It("fetches a new item", func() {
					subject.Invalidate("1")
					_, _ = subject.Get(ignoreCtx, "1")
				})
			})
		})
	})

	When("at capacity", func() {
		BeforeEach(func() {
			subject = cache.NewLRUItem(2, valueMapperWrap(source))
		})

		When("get an uncached item", func() {
			BeforeEach(func() {
				source.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Times(2).Return("1", nil)
				source.EXPECT().Get(gomock.Any(), gomock.Eq("2")).Times(1).Return("2", nil)
				source.EXPECT().Get(gomock.Any(), gomock.Eq("3")).Times(1).Return("3", nil)
				_, _ = subject.Get(ignoreCtx, "1")
				_, _ = subject.Get(ignoreCtx, "2")
			})
			It("removes oldest item", func() {
				Expect(subject.Get(ignoreCtx, "3")).Should(Equal("3"))
				Expect(subject.Get(ignoreCtx, "1")).Should(Equal("1"))
			})
		})

		When("get oldest cached item", func() {
			BeforeEach(func() {
				source.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Times(1).Return("1", nil)
				source.EXPECT().Get(gomock.Any(), gomock.Eq("2")).Times(2).Return("2", nil)
				source.EXPECT().Get(gomock.Any(), gomock.Eq("3")).Times(1).Return("3", nil)
				_, _ = subject.Get(ignoreCtx, "1")
				_, _ = subject.Get(ignoreCtx, "2")
			})
			It("is no longer least recently used", func() {
				// refresh one
				_, _ = subject.Get(ignoreCtx, "1")
				// kick out two
				_, _ = subject.Get(ignoreCtx, "3")
				// reload 2
				_, _ = subject.Get(ignoreCtx, "2")
			})
		})

		When("valueMap fails", func() {
			BeforeEach(func() {
				source.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Times(1).Return("1", nil)
				source.EXPECT().Get(gomock.Any(), gomock.Eq("2")).Times(1).Return("2", nil)
				source.EXPECT().Get(gomock.Any(), gomock.Eq("3")).Times(1).Return("", intentionalErr)
				_, _ = subject.Get(ignoreCtx, "1")
				_, _ = subject.Get(ignoreCtx, "2")
			})
			It("does not evict", func() {
				// fails to evict
				_, _ = subject.Get(ignoreCtx, "3")
				// should not call 1 again as it should still be cached
				_, _ = subject.Get(ignoreCtx, "1")
			})
		})
	})

	When("capacity is zero", func() {
		BeforeEach(func() {
			subject = cache.NewLRUItem(0, valueMapperWrap(source))
			source.EXPECT().Get(gomock.Any(), gomock.Eq("1")).Times(2).Return("1", nil)
			_, _ = subject.Get(ignoreCtx, "1")
		})
		It("does not insert", func() {
			_, err := subject.Get(ignoreCtx, "1")
			Expect(err).Should(MatchError(cache.ErrInsufficientCapacity))
		})
	})

})
