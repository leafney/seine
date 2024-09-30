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
)

type MemoBiz struct {
	memoDao *dao.MemoDao
}

func NewMemoBiz(dao *dao.MemoDao) *MemoBiz {
	return &MemoBiz{
		memoDao: dao,
	}
}

func (b *MemoBiz) Add(ctx context.Context) {

}
