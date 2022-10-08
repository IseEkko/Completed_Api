package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//func main() {
//	logger, _ := zap.NewProduction()
//	defer logger.Sync()
//	url := "hrrps://imooc.com"
//	sugar := logger.Sugar()
//	sugar.Infow("failed to fetch URL",
//		"url", url,
//		"attempt", 3,
//	)
//	sugar.Info()
//}

//Person ..
type Person struct {
	//不能为空并且大于10
	Age      int       `form:"age" binding:"required,gt=10"`
	Name     string    `form:"name" binding:"required"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	r := gin.Default()
	r.POST("/5lmh", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBind(&person); err != nil {
			c.String(500, fmt.Sprint(err))
			return
		}
		c.String(200, fmt.Sprintf("%#v", person))
	})
	r.Run()
}
