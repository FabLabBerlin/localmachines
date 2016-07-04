package mock

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"net/http"
	"net/http/httptest"
)

type Server struct {
	testServer     *httptest.Server
	FbInv          fastbill.Invoice
	testClientId   string
	testCustomerId int64
}

func (s *Server) SetPairClientIdCustomerId(clientId string, customerId int64) {
	s.testClientId = clientId
	s.testCustomerId = customerId
}

func (s *Server) URL() string {
	return s.testServer.URL
}

func NewServer() (s *Server) {
	s = &Server{}
	s.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req fastbill.Request
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		if err := dec.Decode(&req); err != nil {
			panic(err.Error())
		}
		var resp interface{}
		if req.SERVICE == fastbill.SERVICE_CUSTOMER_GET {
			var customers []fastbill.Customer
			if s.testClientId == "" && s.testCustomerId == 0 {
				s.testCustomerId = 123
				customers = []fastbill.Customer{
					{
						CUSTOMER_ID: fmt.Sprintf("%v", s.testCustomerId),
					},
				}
			} else {
				filter := req.FILTER.(map[string]interface{})
				if cn := filter["CUSTOMER_NUMBER"]; cn != s.testClientId {
					customers = []fastbill.Customer{}
				} else {
					customers = []fastbill.Customer{
						{
							CUSTOMER_ID: fmt.Sprintf("%v", s.testCustomerId),
						},
					}
				}
			}
			resp = fastbill.CustomerGetResponse{
				RESPONSE: fastbill.CustomerList{
					Customers: customers,
				},
			}

		} else if req.SERVICE == fastbill.SERVICE_INVOICE_GET {
			resp = fastbill.InvoiceCreateResponse{}
		} else if req.SERVICE == fastbill.SERVICE_INVOICE_CREATE {
			buf, err := json.Marshal(req.DATA)
			if err != nil {
				panic(err.Error())
			}
			if err := json.Unmarshal(buf, &s.FbInv); err != nil {
				panic(err.Error())
			}
			tmp := fastbill.InvoiceCreateResponse{}
			tmp.Response.InvoiceId = 123
			resp = tmp
		} else if req.SERVICE == fastbill.SERVICE_INVOICE_COMPLETE {
			resp = fastbill.InvoiceCompleteResponse{}
		} else if req.SERVICE == fastbill.SERVICE_TEMPLATE_GET {
			r := fastbill.TemplateGetResponse{}
			r.Response.Templates = []fastbill.Template{
				{
					Id:   889700,
					Name: "Makea Standard Template",
				},
			}
			resp = r
		} else {
			panic("unknown service")
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			panic(err.Error())
		}

		w.WriteHeader(http.StatusOK)
	}))
	fastbill.API_URL = s.URL()
	return
}
