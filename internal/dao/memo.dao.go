/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:36
 * @Description:
 */

package dao

import (
	"context"
	"github.com/leafney/seine/pkg/gormx"
)

type MemoDao struct {
	db *gormx.GormDBSvc
}

func NewMemoDao(db *gormx.GormDBSvc) *MemoDao {
	return &MemoDao{
		db: db,
	}
}

func (d *MemoDao) Add(ctx context.Context) {
	d.db.DB.Where("")
}
