package method

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	"gorm.io/gorm"
	"hykoj/help"
	"hykoj/models"
	"log"
	"net/smtp"
	"strconv"
	"time"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户名不能为空",
		})
		return
	}
	if password == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "密码不能为空",
		})
		return
	}
	password = help.Getmd5(password)

	var data models.UserBasic
	err := models.DB.Where("name = ? AND password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "用户名或密码不正确",
			})
			return
		}
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "登录出错",
		})
		log.Println("login error", err)
		return
	}
	token, er := help.Gettoken(&data)
	if er != nil {
		log.Println("gettoken error  ", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "gettoken error",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":  200,
		"token": token,
	})

}

func Register(c *gin.Context) {
	var username = c.PostForm("username")
	var pass1 = c.PostForm("pass1")
	var pass2 = c.PostForm("pass2")
	var mail = c.PostForm("mail")
	var yzm = c.PostForm("yzm")
	tr, _ := models.Redis.Get(models.Ctx, mail).Result()
	if yzm != tr {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}
	if pass1 != pass2 {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "密码不一致",
		})
		return
	}
	password := help.Getmd5(pass1)
	var cnt int64 = 0
	var temp models.UserBasic
	models.DB.Where("name = ?", username).Find(&temp).Count(&cnt)
	if cnt > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户名已存在",
		})
		return
	}
	temp.Name = username
	temp.Password = password
	temp.Identity = help.Getuuid()
	temp.Mail = mail
	models.DB.Create(&temp)
	token, _ := help.Gettoken(&temp)
	c.JSON(200, gin.H{
		"code":  200,
		"token": token,
	})
}

func SendEmail(c *gin.Context) {
	to := c.Query("email")
	var str string = help.Getyzm()
	e := email.NewEmail()
	e.From = "995509224@qq.com"
	e.To = []string{to}
	e.Subject = "hyk online judge注册验证"
	e.Text = []byte("欢迎注册hyk online judge，您的验证码是： ")
	e.HTML = []byte("<h1>" + str + "</h1>")
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "995509224@qq.com", "bohqwgfoektnbcgb", "smtp.qq.com"), &tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	models.Redis.Set(models.Ctx, to, str, time.Second*180)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "发送成功",
	})
}

func Getuserdetail(c *gin.Context) {
	var userid = c.Query("userid")
	if userid == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户id不能为空",
		})
		return
	}
	var temp models.UserBasic
	err := models.GetUserDetail(userid).Find(&temp).Error
	if err != nil {
		log.Println("getuserdetail error, err: ", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "查询出错",
		})
		return
	}
	temp.Password = ""
	c.JSON(200, gin.H{
		"code": 200,
		"data": temp,
	})
}

func Userrank(c *gin.Context) {

	var page, size int64
	var err error
	page, err = strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		fmt.Println("getproblem page strconv error :  ", err)
		log.Panicln("getproblem page strconv error :  ", err)
		return
	}
	size, err = strconv.ParseInt(c.DefaultQuery("size", "20"), 10, 32)
	if err != nil {
		fmt.Println("getproblem page strconv error :  ", err)
		log.Panicln("getproblem page strconv error :  ", err)
		return
	}

	var user []models.UserBasic
	tx := models.DB.Model(&models.UserBasic{}).Omit("password").Order("pass_num DESC, submit_num ASC")
	var cnt int64
	tx.Count(&cnt).Offset(int((page - 1) * size)).Limit(int(size)).Find(&user)
	c.JSON(200, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": cnt,
			"data":  user,
		},
	})

}
