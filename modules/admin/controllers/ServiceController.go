package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/controllers"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
	"time"
)

type ServiceController struct {
	controllers.Controller
}

// 服务列表
func (this *ServiceController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, per_page := utils.GetPageInfo(query)

	count := int64(0)
	re := srbac.Db.Model(&models.Service{}).Count(&count)
	srbac.CheckError(re.Error)

	services := []*models.Service{}
	re = srbac.Db.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&services)
	srbac.CheckError(re.Error)

	this.HTML(c, "./views/admin/service/list.gohtml", map[string]interface{}{
		"menu": "service",
		"title": "服务列表",
		"pager": utils.GetPageHtml(count, page, per_page, query, "/admin/service/list"),
		"services": services,
	})
}

// 添加服务
func (this *ServiceController) Add(c *gin.Context) {
	service := &models.Service{}
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		service = models.NewService(params)
		if service.Validate() && service.Create() {
			c.Redirect(http.StatusFound, "/admin/service/list")
			return
		} else {
			this.SetFailed(c, service.GetError())
		}
	}
	this.HTML(c, "./views/admin/service/add.gohtml", map[string]interface{}{
		"menu": "service",
		"title": "添加服务",
		"service": service,
	})
}

// 编辑服务
func (this *ServiceController) Edit(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/service/list")

	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	service := &models.Service{}
	re := srbac.Db.First(service, id)
	srbac.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		service.SetAttributes(params)
		service.UpdatedAt = time.Now().Unix()

		if service.Validate() && service.Update() {
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, service.GetError())
		}
	}

	this.HTML(c, "./views/admin/service/add.gohtml", map[string]interface{}{
		"menu": "service",
		"title": "编辑服务",
		"service": service,
	})
}

// 删除服务
func (this *ServiceController) Delete(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/course/list", false)

	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	re := srbac.Db.Delete(&models.Service{}, id)
	srbac.CheckError(re.Error)

	this.Redirect(c, referer)
}