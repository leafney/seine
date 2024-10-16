/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 01:36
 * @Description:
 */

package model

import "time"

type Memo struct {
	Id       int64     `json:"id" gorm:"column:id;"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func (c *Memo) TableName() string {
	return "memos"
}
