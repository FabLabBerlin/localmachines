package newsletters

import (
	"bytes"
	"encoding/json"
	"github.com/FabLabBerlin/localmachines/controllers"
	"github.com/astaxie/beego"
	"io"
	"net/http"
	"strings"
)

type Controller struct {
	controllers.Controller
}

type SubscribeRequest struct {
	ApiKey string `json:"apikey"`
	Id     string `json:"id"`
	Email  struct {
		Email string `json:"email"`
	} `json:"email"`
	SendWelcome bool `json:"send_welcome"`
}

// @Title EasylabDev
// @Description Signup to EASY LAB dev newsletters
// @Success 200
// @Failure	401	Not authorized
// @Failure	500	Internal Server Error
// @router /easylab_dev [post]
func (c *Controller) EasylabDev() {
	url := urlPrefix() + "/lists/subscribe.json"
	req := SubscribeRequest{
		ApiKey:      beego.AppConfig.String("mailchimpapikey"),
		Id:          beego.AppConfig.String("mailchimpdevnewsletterlistid"),
		SendWelcome: true,
	}
	req.Email.Email = c.GetString("email")
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(req); err != nil {
		beego.Error("Error marshalling request:", err)
		c.Abort("500")
	}
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		beego.Error("Error subscribing E-Mail:", err)
		c.Abort("500")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buf := bytes.Buffer{}
		io.Copy(&buf, resp.Body)
		if s := buf.String(); !strings.Contains(s, "List_AlreadySubscribed") {
			beego.Error("Status code not 200:", buf.String())
			c.Abort("500")
		}
	}
	c.ServeJSON()
}

func urlPrefix() string {
	apikey := beego.AppConfig.String("mailchimpapikey")
	tmp := strings.Split(apikey, "-")
	if len(tmp) != 2 {
		beego.Error("mailchimpapikey has no part after the dash")
		return ""
	}
	prefix := "https://"
	prefix += tmp[1]
	prefix += ".api.mailchimp.com/2.0/"
	return prefix
}
