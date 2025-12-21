package biz

import (
	"context"
	"errors"

	"github.com/leafney/seine/internal/dal"
	"github.com/leafney/seine/internal/model"
	"github.com/leafney/seine/internal/vmodel"
	"github.com/leafney/seine/pkg/errc"
	"github.com/leafney/seine/pkg/errx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserBiz 用户业务逻辑层
type UserBiz struct {
	userDal *dal.UserDal
}

// NewUserBiz 创建用户业务逻辑层
func NewUserBiz(userDal *dal.UserDal) *UserBiz {
	return &UserBiz{
		userDal: userDal,
	}
}

// Create 创建用户
func (b *UserBiz) Create(ctx context.Context, req *vmodel.UserCreateReq) error {
	// 检查用户名是否已存在
	_, err := b.userDal.FindByUsername(ctx, req.Username)
	if err == nil {
		return errx.ErrorCM(errc.ErrDupUser, "用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errx.ErrorE(err)
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errx.ErrorE(err)
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
		Email:    req.Email,
		Status:   1,
	}

	return b.userDal.Create(ctx, user)
}

// Update 更新用户
func (b *UserBiz) Update(ctx context.Context, req *vmodel.UserUpdateReq) error {
	user, err := b.userDal.FindByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errx.ErrorCM(errc.ErrNotFound, "用户不存在")
		}
		return errx.ErrorE(err)
	}

	// 更新字段
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.Status = req.Status

	return b.userDal.Update(ctx, user)
}

// Delete 删除用户
func (b *UserBiz) Delete(ctx context.Context, id uint) error {
	_, err := b.userDal.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errx.ErrorCM(errc.ErrNotFound, "用户不存在")
		}
		return errx.ErrorE(err)
	}

	return b.userDal.Delete(ctx, id)
}

// Query 查询用户列表
func (b *UserBiz) Query(ctx context.Context, req *vmodel.UserQueryReq) (*vmodel.PageResp, error) {
	users, total, err := b.userDal.QueryWithPage(ctx, req)
	if err != nil {
		return nil, errx.ErrorE(err)
	}

	// 转换为响应格式
	list := make([]*vmodel.UserQueryResp, 0, len(users))
	for _, user := range users {
		list = append(list, &vmodel.UserQueryResp{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			Status:   user.Status,
		})
	}

	return &vmodel.PageResp{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// Login 用户登录
func (b *UserBiz) Login(ctx context.Context, req *vmodel.LoginReq) (*vmodel.LoginResp, error) {
	user, err := b.userDal.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errx.ErrorCM(errc.ErrWrongNamePwd, "用户名或密码错误")
		}
		return nil, errx.ErrorE(err)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errx.ErrorCM(errc.ErrWrongNamePwd, "用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errx.ErrorCM(errc.ErrForbidden, "用户已被禁用")
	}

	// TODO: 生成JWT Token
	// 这里暂时返回示例数据，后续集成 jwtx 库后完善
	return &vmodel.LoginResp{
		AccessToken: "mock-token-" + user.Username,
		ExpiresAt:   0,
	}, nil
}
