package router

import (
	"server/api"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	/*创建了四个子路由组，路径前缀都是 /user：
	  userRouter: 普通认证用户路由
	  userPublicRouter: 公开访问路由
	  userLoginRouter: 登录/注册相关路由（带登录记录中间件）
	  userAdminRouter: 管理员专属路由*/
	userRouter := Router.Group("user")
	userPublicRouter := PublicRouter.Group("user")
	userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())
	userAdminRouter := AdminRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi // 获取用户api实例,通过api实例调用具体的接口方法，可以更方便的使用和封装，避免直接使用具体的实现类。
	{
		userRouter.POST("logout", userApi.Logout)                  //登出
		userRouter.PUT("resetPassword", userApi.UserResetPassword) //重置密码
		userRouter.GET("info", userApi.UserInfo)                   //获取用户信息
		userRouter.PUT("changeInfo", userApi.UserChangeInfo)       //修改用户信息
		userRouter.GET("weather", userApi.UserWeather)             //获取天气信息
		userRouter.GET("chart", userApi.UserChart)                 //获取用户统计信息
	}
	{
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword) //忘记密码
		userPublicRouter.GET("card", userApi.UserCard)                  //获取用户名片
	}
	{
		userLoginRouter.POST("register", userApi.Register) //注册
		userLoginRouter.POST("login", userApi.Login)       //登录
	}
	{
		userAdminRouter.GET("list", userApi.UserList)           //获取用户列表
		userAdminRouter.PUT("freeze", userApi.UserFreeze)       //冻结用户
		userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)   //解冻用户
		userAdminRouter.GET("loginList", userApi.UserLoginList) //获取用户登录列表
	}
}
