package models

type User struct {
	UserID   uint64 `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Email    string `json:"email"`    // 邮箱
	Token    string `gorm:"-"`        // 用户令牌
}

func (u *User) TableName() string {
	return "user"
}
