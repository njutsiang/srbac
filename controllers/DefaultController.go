package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DefaultController struct {
	Controller
}

// 首页
func (this *DefaultController) Index(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/admin/service/list")
}
