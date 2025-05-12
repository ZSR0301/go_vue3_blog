package response

import (
	"server/model/database"

	"github.com/gofrs/uuid"
)

type Login struct {
	User                 database.User `json:"user"` //引用数据库用户模型
	AccessToken          string        `json:"access_token"`
	AccessTokenExpiresAt int64         `json:"access_token_expires_at"`
} //登录响应结构体

type UserCard struct {
	UUID      uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	Address   string    `json:"address"`
	Signature string    `json:"signature"`
} //用户卡片响应结构体

type UserChart struct {
	DateList     []string `json:"date_list"`
	LoginData    []int    `json:"login_data"`
	RegisterData []int    `json:"register_data"`
} //用户图表数据响应结构体
