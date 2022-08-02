package method

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hykoj/help"
	"hykoj/judger"
	"hykoj/models"
	"log"
	"net/http"
	"os"
)

func Submit(c *gin.Context) {
	problemid := c.PostForm("problemidentity")
	userid := c.PostForm("userid")
	lau := c.PostForm("language")
	code := c.PostForm("code")
	path, err := help.CodeSave(code, lau)
	//exec.Command("go", "run", path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Save Error:" + err.Error(),
		})
		return
	}

	sb := &models.SubmitBasic{
		Identity:        help.Getuuid(),
		ProblemIdentity: problemid,
		UserIdentity:    userid,
		Path:            path,
	}

	//编译
	temp := []byte(path)
	for i, _ := range temp {
		if temp[i] == '/' {
			temp[i] = '\\'
		}
	}
	path = string(temp)
	tep, err := os.Getwd()
	path = tep + "\\" + path
	path, err = judger.Complie(lau, path)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   -1,
			"msg":    "编译错误",
			"detail": err.Error(),
		})
		return
	}
	// 代码判断
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemid).Preload("TestCases").First(pb).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Error:" + err.Error(),
		})
		return
	}
	// 提示信息
	var list []map[string]interface{}
	for _, testCase := range pb.TestCases {
		testCase := testCase
		var status int
		var output string
		var sta string
		var runtime, runmem int
		output, status, err, runtime, runmem = judger.Running(lau, path, testCase.Input, testCase.Output, pb.MaxMem, pb.MaxRuntime)
		switch status {
		case -1:
			sta = "Wrong Answer"
		case 1:
			sta = "Accept"
		case 2:
			sta = "Time Limit Exceed"
		case 3:
			sta = "Memory Limit Exceed"
		case 4:
			sta = "Runtime Error"
		}
		list = append(list, map[string]interface{}{
			"status":  sta,
			"output":  output,
			"runtime": runtime,
			"runmem":  runmem,
		})
	}
	if help.Erase(path, lau) != nil {
		log.Println(err)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": list,
	})
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		// 更新 user_basic
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userid).Updates(m).Error
		if err != nil {
			return errors.New("UserBasic Modify Error:" + err.Error())
		}
		// 更新 problem_basic
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemid).Updates(m).Error
		if err != nil {
			return errors.New("ProblemBasic Modify Error:" + err.Error())
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error:" + err.Error(),
		})
		return
	}
}
