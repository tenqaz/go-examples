package main

import (
	"fmt"
	"gin_examples/global"
	"gin_examples/initialize"
	"gin_examples/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type User struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func login(c *gin.Context) {

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非校验错误，其他错误直接返回
			c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": utils.RemoveTopStruct(errs.Translate(global.Trans))})
		return
	}

	if user.UserName != "admin" || user.Password != "123456" {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "you are logged in"})
}

func main() {
	err := initialize.InitTrans("zh")
	if err != nil {
		fmt.Printf("初始化翻译器错误, err = %s", err.Error())
		return
	}

	r := gin.Default()
	r.POST("/login", login)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
