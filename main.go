package main

import (
	"hykoj/models"
	"hykoj/router"
)

func main() {
	models.Init()
	r := router.Getrouter()
	r.Run(":8080")

}
