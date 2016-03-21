package modelTest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FabLabBerlin/localmachines/lib/feedback"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	ZENDESK_EMAIL     = "noreply@example.com"
	ZENDESK_API_TOKEN = "t0pS3cr3t"
)

func TestZenDesk(t *testing.T) {
	Convey("Testing ZenDesk model", t, func() {
		Convey("Testing SubmitTicket", func() {
			Convey("It should invoke an HTTP request as documented in the API docs", func(c C) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						panic(err.Error())
					}
					var data map[string]interface{}
					if err := json.Unmarshal(body, &data); err != nil {
						panic(err.Error())
					}
					data = data["ticket"].(map[string]interface{})
					if s := data["subject"].(string); s != "TestSubject" {
						panic(s)
					}
					requester := data["requester"].(map[string]interface{})
					if s := requester["name"]; s != "Joe" {
						panic(s)
					}
					if s := requester["email"]; s != "joe@example.com" {
						panic(s)
					}
					comment := data["comment"].(map[string]interface{})
					if s := comment["body"]; s != "TestComment" {
						panic(s)
					}
				}))
				defer ts.Close()

				zd := feedback.ZenDesk{
					Email:    ZENDESK_EMAIL,
					ApiToken: ZENDESK_API_TOKEN,
					BaseUrl:  ts.URL,
				}
				ticket := feedback.ZenDeskTicket{
					Requester: feedback.ZenDeskTicketRequester{
						Name:  "Joe",
						Email: "joe@example.com",
					},
					Subject: "TestSubject",
					Comment: feedback.ZenDeskTicketComment{
						Body: "TestComment",
					},
				}
				if err := zd.SubmitTicket(ticket); err != nil {
					panic(err.Error())
				}
			})
		})
	})
}
