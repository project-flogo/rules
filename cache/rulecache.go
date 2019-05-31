package rulecache

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
)

type RedisCacheManager struct {
	client *redis.Client
	cfg    config.CacheConfig
}

//LoadTuples loads tuples of a given tuple type from redis cache into a rulesession
func (rcm *RedisCacheManager) LoadTuples(ctx context.Context, td *model.TupleDescriptor, rs model.RuleSession) error {

	keys, err := rcm.client.SMembers(td.Name).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {

		row, err := rcm.client.HGetAll(key).Result()

		values := make(map[string]interface{})

		for k, v := range row {
			values[k] = v
		}

		if err != nil {
			return err
		}
		//convert values[string]string to values[string]interface{}
		t, err := model.NewTuple(model.TupleType(td.Name), values)

		if err != nil {
			return err
		}

		err = rs.Assert(ctx, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rcm *RedisCacheManager) Init(cfg config.CacheConfig) {
	rcm.cfg.Address = cfg.Address
	rcm.cfg.DB = cfg.DB
	rcm.cfg.Name = cfg.Name
	rcm.cfg.Password = cfg.Password
	rcm.cfg.ServerType = cfg.ServerType

	rcm.client = redis.NewClient(&redis.Options{
		Addr:     rcm.cfg.Address,
		Password: rcm.cfg.Password,
		DB:       rcm.cfg.DB,
	})
}

func (rcm *RedisCacheManager) GetRedisClient() *redis.Client {
	return rcm.client
}
