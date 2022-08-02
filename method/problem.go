package method

import (
	"encoding/json"
	"fmt"
	"hykoj/help"
	"hykoj/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSubmitlist(c *gin.Context) {
	var identity = c.Query("identity")
	var status, err = strconv.Atoi(c.Query("status"))
	if err != nil {
		log.Println("status error, ", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "status 不是数字",
		})
		return
	}
	var userid = c.Query("userid")
	var page, tr = strconv.Atoi(c.Query("page"))
	if tr != nil {
		log.Println("page error, ", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "page 不是数字",
		})
		return
	}
	var size, we = strconv.Atoi(c.Query("size"))
	if we != nil {
		log.Println("size error, ", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "size 不是数字",
		})
		return
	}
	page = size * (page - 1)
	var cnt int64
	var temp []models.SubmitBasic
	models.GetSubmitList(identity, userid, status).Count(&cnt).Offset(page).Limit(size).Find(&temp)
	c.JSON(200, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  temp,
			"count": cnt,
		},
	})
}

func GetProblemDetail(c *gin.Context) {
	var identity = c.Query("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "唯一标识不为空",
		})
		return
	}
	db := models.GetProblemDetail(identity)
	var temp models.ProblemBasic
	var cnt int64
	err := db.Count(&cnt).Find(&temp).Error
	if err != nil {
		log.Println("getproblemdetail error, err: ", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "查询出错",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":  200,
		"count": cnt,
		"data":  temp,
	})
}

func CreatProblem(c *gin.Context) {

	var title = c.PostForm("title")
	var content = c.PostForm("content")
	maxtime, _ := strconv.Atoi(c.PostForm("maxtime"))
	maxmem, _ := strconv.Atoi(c.PostForm("maxmem"))
	category := c.PostFormArray("category")
	testcase := c.PostFormArray("testcase")
	var temp = &models.ProblemBasic{
		Title:      title,
		Content:    content,
		MaxMem:     maxmem,
		MaxRuntime: maxtime,
		Identity:   help.Getuuid(),
	}
	var categories = make([]*models.ProblemCategory, 0)
	for _, id := range category {
		tr, _ := strconv.Atoi(id)
		categories = append(categories, &models.ProblemCategory{
			ProblemId:  temp.ID,
			CategoryId: uint(tr),
		})
	}
	temp.ProblemCategories = categories

	var cases = make([]*models.TestCase, 0)
	for _, cs := range testcase {
		casemap := make(map[string]string)
		json.Unmarshal([]byte(cs), &casemap)

		tr := models.TestCase{
			ProblemIdentity: temp.Identity,
			Input:           casemap["input"],
			Output:          casemap["output"],
		}
		cases = append(cases, &tr)
	}
	temp.TestCases = cases
	models.DB.Create(&temp)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func GetProblemList(c *gin.Context) {
	var page, size int64
	var category string
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
	category = c.Query("category")

	db := models.GetProblemList(c.Query("key"), category)
	tx := make([]models.ProblemBasic, 0)
	var cnt int64 = 0
	page--
	db.Count(&cnt).Offset(int(page * size)).Limit(int(size)).Find(&tx)
	c.JSON(200, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": cnt,
			"data":  tx,
		},
	})
}

func Modifyproblem(c *gin.Context) {

}
