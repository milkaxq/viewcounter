package controller

import (
	"encoding/base64"
	"net/http"
	"os"
	"viewscounter/config"
	"viewscounter/helper"
	"viewscounter/service"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	FilterByIpAndUserAgent(context *gin.Context)
	FindProduct(context *gin.Context)
}

type productController struct {
	productService service.ProductService
}

func NewProductController(productServ service.ProductService) ProductController {
	return &productController{
		productService: productServ,
	}
}

// localhost:8080/goapi/view-counter/?type=
func (c *productController) FilterByIpAndUserAgent(context *gin.Context) {

	config.EnvParser()

	productType := context.Query("type")
	// fmt.Println(productType)

	resRegNoHeaderType := os.Getenv("RES_REG_NO_HEADER")
	resGuidHeaderType := os.Getenv("RES_GUID_HEADER")
	mediaGuidHeaderType := os.Getenv("MEDIA_GUID_HEADER")
	rpAccGuidHeaderType := os.Getenv("RP_ACC_GUID_HEADER")
	rpAccRegNoHeaderType := os.Getenv("RP_ACC_REG_NO_HEADER")

	resGuid := context.GetHeader(resGuidHeaderType)
	eGuid, err := base64.StdEncoding.DecodeString(resGuid)
	resNo := context.GetHeader(resRegNoHeaderType)
	eResNo, err := base64.StdEncoding.DecodeString(resNo)

	media := context.GetHeader(mediaGuidHeaderType)
	eMedia, err := base64.StdEncoding.DecodeString(media)

	rpAccGuid := context.GetHeader(rpAccGuidHeaderType)
	eRpAccGuid, err := base64.StdEncoding.DecodeString(rpAccGuid)
	rpAccNo := context.GetHeader(rpAccRegNoHeaderType)
	eRpAccNo, err := base64.StdEncoding.DecodeString(rpAccNo)

	userAgentLength := len(context.Request.UserAgent())

	switch productType {
	case "resource":
		if err != nil || userAgentLength <= 40 {
			res := helper.BuildErrorResponse("", "Something Went Wrong", helper.EmptyObj{})
			context.JSON(http.StatusBadRequest, res)
		} else {
			c.productService.FilterByResGuid(context.ClientIP(), context.Request.UserAgent(), string(eGuid), string(eResNo))
			res := helper.BuildResponse(true, "OK", helper.EmptyObj{})
			context.JSON(http.StatusOK, res)
		}
	case "media":
		if err != nil || userAgentLength <= 40 {
			res := helper.BuildErrorResponse("", "Something Went Wrong", helper.EmptyObj{})
			context.JSON(http.StatusBadRequest, res)
		} else {
			c.productService.FilterByMediaGuid(context.ClientIP(), context.Request.UserAgent(), string(eMedia))
			res := helper.BuildResponse(true, "OK", helper.EmptyObj{})
			context.JSON(http.StatusOK, res)
		}
	case "rp_acc":
		if err != nil || userAgentLength <= 40 {
			res := helper.BuildErrorResponse("", "Something Went Wrong", helper.EmptyObj{})
			context.JSON(http.StatusBadRequest, res)
		} else {
			c.productService.FilterByRpAccGuid(context.ClientIP(), context.Request.UserAgent(), string(eRpAccGuid), string(eRpAccNo))
			res := helper.BuildResponse(true, "OK", helper.EmptyObj{})
			context.JSON(http.StatusOK, res)
		}
	default:
		res := helper.BuildErrorResponse("Wrong Type", "Something Went Wrong", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
}

//localhost:8080/goapi/find-product/?search=...
func (c *productController) FindProduct(context *gin.Context) {
	search := context.Query("search")
	product, err := c.productService.FindProduct(search)
	if err != nil {
		res := helper.BuildErrorResponse("Failed Search", err.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		if product == nil {
			res := helper.BuildResponse(false, "OK", product)
			context.JSON(http.StatusOK, res)
		} else {
			res := helper.BuildResponse(true, "OK", product)
			context.JSON(http.StatusOK, res)
		}

	}
}
