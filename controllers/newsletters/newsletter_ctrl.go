// /api/newsletters
package newsletters

import (
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/FabLabBerlin/localmachines/lib/mailchimp"
	"github.com/astaxie/beego"
)

type Controller struct {
	controllers.Controller
}

// @Title EasylabDev
// @Description Signup to EASY LAB dev newsletters
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /easylab_dev [post]
func (c *Controller) EasylabDev() {
	req := mailchimp.Subscription{
		ApiKey:      beego.AppConfig.String("mailchimpapikey"),
		Id:          beego.AppConfig.String("mailchimpdevnewsletterlistid"),
		SendWelcome: true,
	}
	req.Email.Email = c.GetString("email")
	if err := req.Submit(); err != nil {
		beego.Error("submit to easylab beta newsletter:", err)
		c.Fail(500)
	}

	c.ServeJSON()
}
