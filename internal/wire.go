/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-09-30 00:58
 * @Description:
 */

package internal

import (
	"github.com/google/wire"
	"github.com/leafney/seine/internal/api"
	"github.com/leafney/seine/internal/biz"
	"github.com/leafney/seine/internal/dao"
)

var Set = wire.NewSet(
	api.NewMemoApi,
	biz.NewMemoBiz,
	dao.NewMemoDao,
)
