package middle

import (
	"github.com/gin-gonic/gin"
	"hykoj/help"
)

func Isroot() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		var temp, _ = help.Parsetoken(token)
		if temp == nil || temp.Isroot != 1 {
			c.JSON(200, gin.H{
				"code":  -1,
				"msg":   "不是管理员",
				"temp":  temp,
				"token": token,
			})
			c.Abort()
		}
		c.Next()
	}
}
