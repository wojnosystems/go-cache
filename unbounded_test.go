package cache_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wojnosystems/go-cache"
)

func valueMapperWrap(source *MockstatefulValueMapper) cache.ValueMapper {
	return func(ctx context.Context, key interface{}) (value interface{}, err error) {
		return source.Get(ctx, key.(string))
	}
}

var intentionalErr = fmt.Errorf("intentional")

const (
	expectedValue = "echo"
)

var ignoreCtx = context.Background()

var _ = Describe("StringKey", func() {
	var (
		ctrl   *gomock.Controller
		source *MockstatefulValueMapper
		cacher cache.GetInvalidater
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(ginkgoTestReporter{})
		source = NewMockstatefulValueMapper(ctrl)
		cacher = cache.NewUnbounded(valueMapperWrap(source))
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	When("value factory always succeeds", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, expectedValue).Return(expectedValue, nil)
		})

		When("empty", func() {
			It("returns the item", func() {
				Expect(cacher.Get(ignoreCtx, expectedValue)).Should(Equal(expectedValue))
			})
		})
	})

	When("value already cached", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, "existing").
				Times(1).
				Return("value", nil)
		})
		It("does not call supply", func() {
			Expect(cacher.Get(ignoreCtx, "existing")).Should(Equal("value"))
			Expect(cacher.Get(ignoreCtx, "existing")).Should(Equal("value"))
		})
	})

	When("value factory always fails", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, "some key").Return("", intentionalErr)
		})

		When("get is called", func() {
			It("fails", func() {
				_, err := cacher.Get(ignoreCtx, "some key")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	When("invalidated", func() {
		BeforeEach(func() {
			source.EXPECT().Get(ignoreCtx, "1").Times(2).Return("1", nil)
			_, _ = cacher.Get(ignoreCtx, "1")
		})

		It("fetches a new item", func() {
			cacher.Invalidate("1")
			_, _ = cacher.Get(ignoreCtx, "1")
		})
	})
})
