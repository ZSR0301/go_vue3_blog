package request

type Register struct {
	Username         string `json:"username" binding:"required,max=20"` //binding定义验证规则
	Password         string `json:"password" binding:"required,min=8,max=16"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
} //用户注册接口的请求体

type Login struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=16"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
} //用户登录接口的请求体

type ForgotPassword struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
	NewPassword      string `json:"new_password" binding:"required,min=8,max=16"`
} //忘记密码重置接口的请求体

type UserCard struct {
	UUID string `json:"uuid" form:"uuid" binding:"required"`
} //获取用户卡片信息的请求参数

type UserResetPassword struct {
	UserID      uint   `json:"-"` //忽略 UserID 字段，因为它是由 JWT 解析出来的，并不在请求体中，防止前端越权，直接通过后端读取。
	Password    string `json:"password" binding:"required,min=8,max=16"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=16"`
} //登录状态下修改密码

type UserChangeInfo struct {
	UserID    uint   `json:"-"`
	Username  string `json:"username" binding:"required,max=20"`
	Address   string `json:"address" binding:"max=200"`
	Signature string `json:"signature" binding:"max=320"`
} //修改用户个人信息

type UserChart struct {
	Date int `json:"date" form:"date" binding:"required,oneof=7 30 90 180 365"`
} //获取用户图表数据的时间范围参数

type UserList struct {
	UUID     *string `json:"uuid" form:"uuid"`
	Register *string `json:"register" form:"register"`
	PageInfo
} //管理员获取用户列表的查询条件

type UserOperation struct {
	ID uint `json:"id" binding:"required"`
} //操作用户的通用请求(如冻结/解冻)

type UserLoginList struct {
	UUID *string `json:"uuid" form:"uuid"`
	PageInfo
} //查询用户登录记录
