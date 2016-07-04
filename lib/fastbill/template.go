package fastbill

import (
	"fmt"
)

type Template struct {
	Id   int64  `json:"TEMPLATE_ID,omitempty"`
	Name string `json:"TEMPLATE_NAME,omitempty"`
}

type TemplateGetResponse struct {
	Request  Request `json:"REQUEST,omitempty"`
	Response struct {
		Templates []Template `json:"TEMPLATES,omitempty"`
	} `json:"RESPONSE,omitempty"`
}

func ListTemplates() (ts []Template, err error) {
	fb := New()
	request := Request{
		SERVICE: SERVICE_TEMPLATE_GET,
	}
	var response TemplateGetResponse
	if err = fb.execGetRequest(&request, &response); err != nil {
		return nil, fmt.Errorf("get request: %v", err)
	}
	return response.Response.Templates, err
}

func TemplateIdExists(id int64) (yes bool, err error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid id: %v", id)
	}

	ts, err := ListTemplates()
	if err != nil {
		return false, fmt.Errorf("list templates: %v", err)
	}
	for _, t := range ts {
		if t.Id == id {
			return true, nil
		}
	}
	return false, nil
}
