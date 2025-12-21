package dal

import (
	"context"

	"github.com/leafney/seine/internal/model"
	"github.com/leafney/seine/internal/vmodel"
	"gorm.io/gorm"
)

// UserDal 用户数据访问层
type UserDal struct {
	db *gorm.DB
}

// NewUserDal 创建用户数据访问层
func NewUserDal(db *gorm.DB) *UserDal {
	return &UserDal{db: db}
}

// Create 创建用户
func (d *UserDal) Create(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据ID查找用户
func (d *UserDal) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := d.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (d *UserDal) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (d *UserDal) Update(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (d *UserDal) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// QueryWithPage 分页查询用户
func (d *UserDal) QueryWithPage(ctx context.Context, req *vmodel.UserQueryReq) ([]*model.User, int64, error) {
	var users []*model.User
	var count int64

	query := d.db.WithContext(ctx).Model(&model.User{})

	// 构建查询条件
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 获取总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error

	return users, count, err
}
