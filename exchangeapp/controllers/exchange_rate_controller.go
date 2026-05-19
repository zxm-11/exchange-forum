package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateExchangeRate函数用于创建汇率
func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate
	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	exchangeRate.Date = time.Now() //设置当前时间为汇率的日期
	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, exchangeRate)
}

// GetExchangeRate函数用于获取汇率
func GetExchangeRate(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate
	//find方法用于查询数据库中的数据,参数是一个指向切片的指针,查询结果会被存储在这个切片中
	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, exchangeRates)
}
