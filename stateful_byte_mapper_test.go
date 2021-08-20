//go:generate mockgen -source=stateful_byte_mapper_test.go -destination=stateful_byte_mapper_mock_test.go -package=cache_test
package cache

import "context"

type statefulByteMapper interface {
	Get(ctx context.Context, key string) (value []byte, err error)
}
