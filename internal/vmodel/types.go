package vmodel

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,max=64"`
	Name     string `json:"name" validate:"max=64"`
	Email    string `json:"email" validate:"email,max=128"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	ID       uint   `json:"id" validate:"required"`
	Name     string `json:"name" validate:"max=64"`
	Email    string `json:"email" validate:"email,max=128"`
	Status   int    `json:"status" validate:"oneof=0 1"`
}

// UserQueryReq 查询用户请求
type UserQueryReq struct {
	Page     int    `json:"page" validate:"required,min=1"`
	PageSize int    `json:"pageSize" validate:"required,min=1,max=100"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// UserQueryResp 查询用户响应
type UserQueryResp struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResp 登录响应
type LoginResp struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// PageResp 分页响应
type PageResp struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
