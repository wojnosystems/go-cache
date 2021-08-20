package cache_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-cache"
)

var _ = Describe("LruByte", func() {
	var (
		ctrl     *gomock.Controller
		source   *MockstatefulByteMapper
		capacity uint
		subject  cache.ByteGetInvalidator
	)
	BeforeEach(func() {
		capacity = 100
		ctrl = gomock.NewController(ginkgoTestReporter{})
		source = NewMockstatefulByteMapper(ctrl)
		subject = cache.NewLRUByte(capacity, func(ctx context.Context, key interface{}) (value []byte, err error) {
			v, err := source.Get(ctx, key.(string))
			return v, err
		})
	})
	AfterEach(func() {
		ctrl.Finish()
	})
	When("cannot fit", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, "100 + 10").Times(1).
				Return(make([]byte, 100+10), nil)
		})
		It("does not cache", func() {
			_, err := subject.Get(ignoreCtx, "100 + 10")
			Expect(err).Should(HaveOccurred())
		})
	})
	When("items fit", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, "60").Times(2).
				Return(make([]byte, 60), nil)
			source.EXPECT().Get(ignoreCtx, "50").Times(2).
				Return(make([]byte, 50), nil)
		})
		It("evicts items", func() {
			_, err := subject.Get(ignoreCtx, "60")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = subject.Get(ignoreCtx, "50")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = subject.Get(ignoreCtx, "60")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = subject.Get(ignoreCtx, "50")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	When("invalidated", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, "20").Times(2).
				Return(make([]byte, 20), nil)
		})
		It("forces a fetch", func() {
			_, _ = subject.Get(ignoreCtx, "20")
			subject.Invalidate("20")
			_, _ = subject.Get(ignoreCtx, "20")
		})
	})
})
