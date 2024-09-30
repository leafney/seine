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
	"gorm.io/gorm"
)

type MemoDao struct {
	db *gorm.DB
}

func NewMemoDao(db *gormx.DBService) *MemoDao {
	return &MemoDao{
		db: db.DB,
	}
}

func (d *MemoDao) Add(ctx context.Context) {

}
