/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:35
 * @Description:
 */

package biz

import (
	"context"
	"github.com/leafney/seine/internal/dao"
	"github.com/leafney/seine/pkg/redisx"
)

type MemoBiz struct {
	cache   *redisx.RRedisSvc
	memoDao *dao.MemoDao
}

func NewMemoBiz(dao *dao.MemoDao, cache *redisx.RRedisSvc) *MemoBiz {
	return &MemoBiz{
		memoDao: dao,
		cache:   cache,
	}
}

func (b *MemoBiz) Add(ctx context.Context) {

}
