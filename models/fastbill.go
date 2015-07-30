package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

// API endpoint - all requests go here
const (
	FASTBILL_API_URL                 = "https://my.fastbill.com/api/1.0/api.php"
	FASTBILL_SERVICE_CUSTOMER_GET    = "customer.get"
	FASTBILL_SERVICE_CUSTOMER_CREATE = "customer.create"
	FASTBILL_CUSTOMER_TYPE_BUSINESS  = "business"
	FASTBILL_CUSTOMER_TYPE_CONSUMER  = "consumer"
)

// Main FastBill object. All the functionality goes through this object.
type FastBill struct {
	Email  string
	APIKey string
}

// Base FastBill API request model
type FastBillRequest struct {
	LIMIT    int64
	OFFSET   int64
	SERVICE  string
	REQUEST  interface{}
	FILTER   interface{}
	DATA     interface{}
	RESPONSE interface{}
	ERRORS   interface{}
}

// customer.get response model
// For now define a response per request. The interface{} thing
// does not work well with the JSON unmarshal thing for some reason.
// But this is more clear in a way.
type FastBillCustomerGetResponse struct {
	REQUEST  FastBillRequest
	RESPONSE FastBillCustomerList
}

// Response model that we expect from the FastBill API on customer.create
// request
type FastBillCustomerCreateResponse struct {
	REQUEST  FastBillRequest
	RESPONSE struct {
		STATUS      string
		CUSTOMER_ID int64
	}
}

// Response model for the JSON that will be returned to this API clients
type FastBillCreateCustomerResponse struct {
	CUSTOMER_ID int64
}

// FastBill customer model
type FastBillCustomer struct {
	CUSTOMER_ID                    string
	CUSTOMER_NUMBER                string
	DAYS_FOR_PAYMENT               string
	CREATED                        string
	PAYMENT_TYPE                   string
	BANK_NAME                      string
	BANK_ACCOUNT_NUMBER            string
	BANK_CODE                      string
	BANK_ACCOUNT_OWNER             string
	BANK_IBAN                      string
	BANK_BIC                       string
	BANK_ACCOUNT_MANDATE_REFERENCE string
	SHOW_PAYMENT_NOTICE            string
	ACCOUNT_RECEIVABLE             string
	CUSTOMER_TYPE                  string // Required. Customer type: business | consumer
	TOP                            string
	NEWSLETTER_OPTIN               string
	CONTACT_ID                     string
	ORGANIZATION                   string // Company name [REQUIRED] when CUSTOMER_TYPE = business
	POSITION                       string
	SALUTATION                     string
	FIRST_NAME                     string
	LAST_NAME                      string // Last name [REQUIRED] when CUSTOMER_TYPE = consumer
	ADDRESS                        string
	ADDRESS_2                      string
	ZIPCODE                        string
	CITY                           string
	COUNTRY_CODE                   string
	SECONDARY_ADDRESS              string
	PHONE                          string
	PHONE_2                        string
	FAX                            string
	MOBILE                         string
	EMAIL                          string
	VAT_ID                         string
	CURRENCY_CODE                  string
	LASTUPDATE                     string
	TAGS                           string
}

// Customer list model
type FastBillCustomerList struct {
	Customers []FastBillCustomer
}

// Filter model for customer.get request
type FastBillCustomerGetFilter struct {
	CUSTOMER_ID     string
	CUSTOMER_NUMBER string
	COUNTRY_CODE    string
	CITY            string
	TERM            string // Search term in one of the given fields: ORGANIZATION, FIRST_NAME, LAST_NAME, ADDRESS, ADDRESS_2, ZIPCODE, EMAIL, TAGS.
}

// Get customers with support for limit and offset
func (this *FastBill) GetCustomers(filter *FastBillCustomerGetFilter,
	limit int64, offset int64) (*FastBillCustomerList, error) {

	request := FastBillRequest{}
	request.SERVICE = FASTBILL_SERVICE_CUSTOMER_GET
	request.LIMIT = limit
	request.OFFSET = offset
	request.FILTER = filter

	response := FastBillCustomerGetResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute get request: %v", err)
	}

	return &response.RESPONSE, nil
}

// Create FastBill customer, returns Customer ID
func (this *FastBill) CreateCustomer(customer *FastBillCustomer) (int64, error) {

	request := FastBillRequest{}
	request.SERVICE = FASTBILL_SERVICE_CUSTOMER_CREATE
	request.DATA = customer

	response := FastBillCustomerCreateResponse{}
	err := this.execGetRequest(&request, &response)
	if err != nil {
		return 0, fmt.Errorf("Failed to execute get request: %v", err)
	}

	if response.RESPONSE.STATUS == "success" {
		return response.RESPONSE.CUSTOMER_ID, nil
	} else {
		return 0, errors.New("There was an error while creating a customer")
	}
}

// Reusable helper function for the
func (this *FastBill) execGetRequest(request *FastBillRequest, response interface{}) error {

	var err error
	var req *http.Request
	var resp *http.Response
	var jsonBytes []byte

	jsonBytes, err = json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}
	beego.Trace(string(jsonBytes))

	req, err = http.NewRequest("GET", FASTBILL_API_URL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(this.Email, this.APIKey)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal json: %v", err)
	}

	return nil
}

/* Country codes used by FastBill API

14    Afghanistan
13    Ägypten
12    Albanien
3     Algerien
4     Andorra
5     Angola
6     Antigua und Barbuda
193   Äquatorialguinea
7     Argentinien
8     Armenien
9     Aserbeidschan
194   Äthiopien
10    Australien
15    Bahamas
16    Bahrain
22    Bangladesh
18    Barbados
17    Belgien
19    Belize
23    Benin
20    Bhutan
21    Bolivien
212   Bonaire
24    Bosnien
25    Botsuana
26    Brasilien
200   Britische Jungferninseln
27    Brunei
28    Bulgarien
29    Burkina Faso
30    Burundi
31    Chile
32    China
33    Costa Rica
34    Cote d´Ivoire
38    Dänemark
1     elected>Deutschland
210   Dominica
36    Dominikanische Republik
37    Dschibuti
39    Ecuador
40    El Salvador
2     England
41    Eritrea
42    Estland
43    Finnland
44    Frankreich
216   Französisch-Polynesien
45    Gabun
46    Gambia
47    Georgien
48    Ghana
208   Gibraltar
49    Grenada
50    Griechenland
51    Großbritannien
52    Guatemala
53    Guinea
54    Guinea-Bissau
55    Guyana
56    Haiti
57    Honduras
202   Hong Kong
58    Indien
59    Indonesien
203   Insel Man
60    Irak
61    Iran
198   Irland
62    Island
63    Israel
64    Italien
65    Jamaika
66    Japan
67    Jemen
68    Jordanien
69    Jugoslawien
204   Kaimaninseln
70    Kambodscha
71    Kamerun
72    Kanada
73    Kap Verde
74    Kasachstan
75    Katar
76    Kenia
77    Kirgisistan
78    Kiribati
79    Kolumbien
80    Komoren
81    Kongo
205   Kosovo
82    Kroatien
83    Kuba
84    Kuweit
85    Laos
86    Laos
87    Lesotho
88    Lettland
89    Libanon
90    Liberia
91    Libyen
92    Liechtenstein
93    Litauen
94    Luxemburg
95    Madagaskar
96    Malawi
97    Malaysia
98    Malediven
99    Mali
100   Malta
101   Marokko
102   Marshall-Inseln
201   Martinique
103   Mauretanien
104   Mauritius
105   Mazedonien
106   Mexiko
107   Mikronesien
108   Moldavien
109   Moldavien
110   Monaco
111   Mongolei
112   Montenegro
113   Mosambik
114   Myanmar
115   Namibia
116   Nauru
117   Nepal
118   Neuseeland
119   Nicaragua
120   Niederlande
214   Niederländische Antillen
121   Niger
122   Nigeria
123   Nordkorea
124   Norwegen
11    Oman
195   Österreich
125   Pakistan
127   Palästina
126   Palau
128   Panama
129   Papua Neuguinea
130   Paraguay
131   Peru
197   Philippinen
133   Polen
134   Portugal
136   Ruanda
137   Rumänien
138   Russland
143   Salomonen
144   Sambia
145   Sambia
146   Samoa
147   San Marino
148   Saudi-Arabien
149   Schweden
150   Schweiz
151   Senegal
152   Serbien
153   Seychellen
154   Sierra Leone
155   Sierra Leone
156   Simbabwe
157   Singapur
158   Slowakei
159   Slowenien
160   Somalia
161   Spanien
162   Sri Lanka
140   St. Kitts u. Nevis
141   St. Lucia
142   St. Vincent
135   Südafrika
163   Sudan
166   Südkorea
207   Südsudan
164   Suriname
165   Swasiland
167   Syrien
168   Tadschikistan
169   Taiwan
170   Tansania
171   Thailand
172   Togo
173   Tonga
174   Trinidad und Tobago
175   Tschad
176   Tschechische Republik
177   Tunesien
180   Türkei
178   Turkmenistan
179   Tuvalu
181   Uganda
182   Ukraine
196   Ungarn
183   Uruguay
184   USA
185   Usbekistan
186   Vanuatu
187   Vatikanstadt
188   Venezuela
189   Vereinigte Arab. Emirate
190   Vietnam
206   Vojvodina
191   Weissrussland
139   Westsahara
192   Zentralafrikanische Republik
199   Zypern

*/
