package middleware

import (
	"server/model/appTypes"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// AdminAuth 是一个中间件，用于检查用户是否具有管理员权限
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户的角色ID
		roleID := utils.GetRoleID(c)

		// 检查用户是否为管理员
		if roleID != appTypes.Admin {
			// 如果不是管理员，返回访问被拒绝的响应
			response.Forbidden("Access denied. Admin privileges are required", c)
			// 中止请求处理
			c.Abort()
			return
		}

		// 如果用户是管理员，继续执行后续处理器
		c.Next()
	}
}
