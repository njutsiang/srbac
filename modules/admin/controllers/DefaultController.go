package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/controllers"
)

type DefaultController struct {
	controllers.Controller
}

// 后台首页
func (this *DefaultController) Index(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/admin/service/list")
}