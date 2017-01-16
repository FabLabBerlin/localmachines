package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"net/http"
	"strings"
)

type Subscription struct {
	ApiKey string `json:"apikey"`
	Id     string `json:"id"`
	Email  struct {
		Email string `json:"email"`
	} `json:"email"`
	SendWelcome bool `json:"send_welcome"`
}

func (req Subscription) Submit() (err error) {
	url := urlPrefix(req.ApiKey) + "/lists/subscribe.json"
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(req); err != nil {
		return fmt.Errorf("Error marshalling request: %v", err)
	}
	resp, err := http.Post(url, "application/json", &buf)
	if err != nil {
		return fmt.Errorf("Error subscribing E-Mail: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		buf := bytes.Buffer{}
		io.Copy(&buf, resp.Body)
		if s := buf.String(); !strings.Contains(s, "List_AlreadySubscribed") {
			return fmt.Errorf("Status code not 200: %v", buf.String())
		}
	}

	return
}

func urlPrefix(apikey string) string {
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
