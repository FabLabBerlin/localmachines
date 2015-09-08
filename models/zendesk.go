package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type ZenDesk struct {
	Email    string
	ApiToken string
	BaseUrl  string // https://<subdomain>.zendesk.com
}

func (this ZenDesk) ApiUrl(function string) string {
	return this.BaseUrl + "/api/v2/" + function
}

func (this ZenDesk) SubmitTicket(ticket ZenDeskTicket) error {
	request := ZenDeskTicketRequest{
		Ticket: ticket,
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}
	url := this.ApiUrl("tickets.json")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(this.Email+"/token", this.ApiToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error("Failed to get response:", err)
		return fmt.Errorf("Failed to get response")
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("Failed to read response body:", err)
		return fmt.Errorf("Failed to read response body")
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("unexpected Status Code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}

type ZenDeskTicketRequest struct {
	Ticket ZenDeskTicket `json:"ticket"`
}

type ZenDeskTicket struct {
	Requester ZenDeskTicketRequester `json:"requester"`
	Subject   string                 `json:"subject"`
	Comment   ZenDeskTicketComment   `json:"comment"`
}

type ZenDeskTicketRequester struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ZenDeskTicketComment struct {
	Body string `json:"body"`
}
