package mocks

import (
	"encoding/json"
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"os"
)

func LoadInvoice(id int64) (inv *invutil.Invoice) {
	if err := load("invoices", id, &inv); err != nil {
		panic(err.Error())
	}
	return
}

func LoadUserMembership(id int64) (um *user_memberships.UserMembership) {
	if err := load("user_memberships", id, &um); err != nil {
		panic(err.Error())
	}
	return
}

func load(model string, id int64, v interface{}) (err error) {
	fn := fmt.Sprintf("./models/mocks/%v/%v.json", model, id)

	f, err := os.Open(fn)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(v); err != nil {
		panic(err.Error())
	}

	return
}
