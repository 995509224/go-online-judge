package method

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hykoj/help"
	"hykoj/models"
	"log"
	"strconv"
)

func Getcategorylist(c *gin.Context) {
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
	var cnt int64
	var data []models.CategoryBasic
	models.DB.Model(&models.CategoryBasic{}).Find(&data).Offset(int((page - 1) * size)).Limit(int(size)).Count(&cnt)
	c.JSON(200, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": cnt,
			"data":  data,
		},
	})

}

func CreatCategory(c *gin.Context) {
	var name = c.PostForm("name")
	var parentid, _ = strconv.Atoi(c.DefaultPostForm("parentid", "0"))

	if name == "" {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "参数不正确",
		})
		return
	}

	var temp = models.CategoryBasic{
		Name:     name,
		ParentId: parentid,
		Identity: help.Getuuid(),
	}

	var cnt int64

	models.DB.Model(&models.CategoryBasic{}).Where("name = ?", name).Count(&cnt)
	if cnt > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "已存在该分类",
		})
		return
	}
	err := models.DB.Create(&temp).Error
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "创建错误",
		})
	}
}

func Modifycategory(c *gin.Context) {
	var name = c.PostForm("name")
	var parentid, _ = strconv.Atoi(c.DefaultPostForm("parentid", "0"))
	var uuid = c.PostForm("identity")

	var temp = models.CategoryBasic{
		Name:     name,
		ParentId: parentid,
		Identity: uuid,
	}
	err := models.DB.Model(&models.CategoryBasic{}).Where("identity = ?", uuid).Updates(temp).Error
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "删除成功",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

func DeleteCategory(c *gin.Context) {
	var name = c.PostForm("identity")
	if name == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	err := models.DB.Model(&models.ProblemCategory{}).Where("category_id in (select id from category_basic where name = ?)", name).Delete(&models.ProblemBasic{})
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "删除错误",
		})

		return
	}
	err = models.DB.Model(&models.CategoryBasic{}).Where("name = ?", name).Delete(&models.CategoryBasic{})
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "删除错误",
		})

		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
