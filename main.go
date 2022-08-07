package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

func main() {
	dsn := "go-crud:g0-crud@tcp(qiuqian.xyz:6603)/go-crud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	fmt.Println(db, err)

	sqlDB, err := db.DB()

	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); notnull" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20); notnull" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); notnull" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); notnull" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200); notnull" json:"address" binding:"required"`
	}

	db.AutoMigrate(&List{})
	r := gin.Default()

	// 增
	r.POST("/user/add", func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data)

		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			db.Create(&data)
			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})

	// 删
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List
		id := c.Param("id")

		db.Where("id = ?", id).Find(&data)

		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id未找到，删除失败",
				"code": 400,
			})
		} else {
			db.Where("id = ?", id).Delete(&data)
			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 改
	r.PUT("/user/update/:id", func(c *gin.Context) {
		var data List
		id := c.Param("id")

		db.Select("id").Where("id = ?", id).Find(&data)

		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户ID不存在",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				db.Where("id = ?", id).Updates(&data)
				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})

	// 查
	r.GET("/user/list", func(c *gin.Context) {
		var dataList []List
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}

		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		var total int64
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "查询成功",
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
				"code": 200,
			})
		}
	})
	r.GET("/user/list/:name", func(c *gin.Context) {
		name := c.Param("name")

		var dataList []List
		db.Where("name = ?", name).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到" + name + "的数据",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"data": dataList,
				"code": 200,
			})
		}
	})

	PORT := "8080"
	r.Run(":" + PORT)
}
