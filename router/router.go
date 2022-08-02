package router

import (
	"github.com/gin-gonic/gin"
	"hykoj/method"
	"hykoj/middle"
)

func Getrouter() *gin.Engine {

	h := gin.Default()
	h.GET("/problemlist", method.GetProblemList)
	h.GET("/problemdetail", method.GetProblemDetail)
	h.GET("/userdetil", method.Getuserdetail)
	h.GET("/submitlist", method.GetSubmitlist)
	h.POST("/login", method.Login)
	h.GET("/sendemail", method.SendEmail)
	h.POST("/register", method.Register)
	h.GET("/userrank", method.Userrank)
	h.GET("/categorylist", method.Getcategorylist)
	h.POST("/submit", method.Submit)

	h.POST("/creatproblem", middle.Isroot(), method.CreatProblem)
	h.POST("/creatcategory", middle.Isroot(), method.CreatCategory)
	h.PUT("/modifycategory", middle.Isroot(), method.Modifycategory)
	h.DELETE("/deletecategory", middle.Isroot(), method.DeleteCategory)
	//h.PUT("/modifyproblem", middle.Isroot(), method.Modifyproblem)
	return h

}
