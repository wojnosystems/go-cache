//go:generate mockgen -source=stateful_value_mapper_test.go -destination=stateful_value_mapper_mock_test.go -package=cache_test
package cache

import "context"

type statefulValueMapper interface {
	Get(ctx context.Context, key string) (value string, err error)
}
