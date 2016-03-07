package fastbill

import (
	"encoding/json"
	"testing"

	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	. "github.com/smartystreets/goconvey/convey"
)

const FASTBILL_RESPONSE = `
{
  "REQUEST":{
    "SERVICE":"invoice.create",
    "DATA":{
      "CUSTOMER_ID":"696",
      "TEMPLATE_ID":"1",
      "ITEMS":[
        {
          "DESCRIPTION":"3D Print Club",
          "QUANTITY":"1",
          "UNIT_PRICE":"10",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 4 Mia (Replicator 5 gen)",
          "QUANTITY":"0.23803215168333333",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 5 Pumpkin (I3 Berlin)",
          "QUANTITY":"385.42372557408333",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 6 Honey Bunny (I3 Berlin)",
          "QUANTITY":"1398.8801021362833",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 7 Fabienne (i3 Berlin)",
          "QUANTITY":"0.34533285080000004",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - Replicator Mini",
          "QUANTITY":"2985.365968979117",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"Electronics Desk",
          "QUANTITY":"106.27382089261667",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"Laser Cutter - Epilog Zing 6030",
          "QUANTITY":"0.8469698362333333",
          "UNIT_PRICE":"0.8",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 1 Vincent Vega (Replicator 2)",
          "QUANTITY":"0.8576845149000001",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"Vinyl Cutter",
          "QUANTITY":"19.81636890671667",
          "UNIT_PRICE":"0.5",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 9 Angelo (Replicator 2)",
          "QUANTITY":"2741.395774145317",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        },
        {
          "DESCRIPTION":"3D Printer - 2 Jules (Replicator 5 gen)",
          "QUANTITY":"858.6029838759333",
          "UNIT_PRICE":"0.1",
          "VAT_PERCENT":"19",
          "SORT_ORDER":""
        }
      ]
    }
  },
  "RESPONSE":{
    "STATUS":"success",
    "INVOICE_ID":4833738,
    "DOCUMENT_ID":""
  }
}
`

func TestFastbillInvoice(t *testing.T) {
	Convey("Testing Fastbill Invoice", t, func() {
		Convey("Testing whether Fastbill responses can be unmarshaled", func() {
			var response fastbill.InvoiceCreateResponse
			err := json.Unmarshal([]byte(FASTBILL_RESPONSE), &response)
			So(err, ShouldBeNil)
		})
	})
}
