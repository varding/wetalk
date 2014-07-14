package post

import (
	"github.com/beego/wetalk/modules/models"
	"github.com/beego/wetalk/routers/base"
)

type NoticeRouter struct {
	base.BaseRouter
}

func (this *NoticeRouter) Get() {
	this.Data["IsHomePage"] = true
	this.TplNames = "post/notice.html"

	if this.CheckLoginRedirect() {
		return
	}

	var notifications []models.Notification
	qs := models.Notifications(this.User.Id)

	pers := 10
	count, _ := models.CountObjects(qs)
	pager := this.SetPaginator(pers, count)

	qs = qs.OrderBy("-Created").Limit(pers, pager.Offset()).RelatedSel()

	models.ListObjects(qs, &notifications)
	this.Data["Notifications"] = notifications
}
