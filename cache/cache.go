package rulecache

import (
	"context"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
)

type CacheManager interface {
	Init(cfg config.CacheConfig)
	LoadTuples(ctx context.Context, td *model.TupleDescriptor, rs model.RuleSession) error
}
