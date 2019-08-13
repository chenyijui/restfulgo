package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"restfulgo/model"
	_ "restfulgo/model"
	"time"
)
var (
	db *gorm.DB
	sqlConnection = "root:@/test?charset=utf8&parseTime=True&loc=Local"
)

func init()  {
	router := gin.Default()
	router.GET("/user/", getUsers)
	router.GET("/user/:id", getOneUser)
	router.POST("/user", createUser)
	router.PUT("/user/:id", updateUser)
	router.DELETE("/user/:id", deleteUser)
}

func main() {
	db, err := gorm.Open("mysql", sqlConnection)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})

	router.Run(":6999")

}
func getUsers(c *gin.Context) {
	users := make([]model.User, 0)
	//db.Find(&users)
	db.Scopes(AgeGreaterThan20).Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func AgeGreaterThan20(db *gorm.DB) *gorm.DB {
	return db.Where("age > ?", 20)
}

func getOneUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}

func createUser(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)
	user.Birthday = time.Now()
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.BindJSON(&user)
		db.Save(&user)
		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	}
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	db.Where("id = ?", id).Delete(&user)
	c.JSON(http.StatusOK, gin.H{
		"data": "this has been deleted!",
	})
}
