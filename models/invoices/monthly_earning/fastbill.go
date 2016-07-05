package monthly_earning

import (
	"github.com/FabLabBerlin/localmachines/lib/fastbill"
	"github.com/FabLabBerlin/localmachines/models/user_roles"
	"github.com/astaxie/beego"
)

type DraftsCreationReport struct {
	Ids                 []int64
	SuccessUids         []int64
	EmptyUids           []int64
	AlreadyExportedUids []int64
	Errors              []DraftsCreationError
}

type DraftsCreationError struct {
	UserId  int64
	Problem string
}

func CreateFastbillDrafts(me *MonthlyEarning, vatPercent float64) (report DraftsCreationReport) {
	report.Ids = make([]int64, 0, len(me.Invoices))
	report.SuccessUids = make([]int64, 0, len(me.Invoices))
	report.EmptyUids = make([]int64, 0, len(me.Invoices))
	report.AlreadyExportedUids = make([]int64, 0, len(me.Invoices))
	report.Errors = make([]DraftsCreationError, 0, len(me.Invoices))

	for _, inv := range me.Invoices {
		uid := inv.User.Id
		inv.VatPercent = vatPercent
		if inv.User.NoAutoInvoicing {
			continue
		}
		if r := inv.User.GetRole(); (r == user_roles.STAFF || r == user_roles.ADMIN || r == user_roles.SUPER_ADMIN) && uid != 19 {
			e := DraftsCreationError{
				UserId:  uid,
				Problem: "User role is " + r.String(),
			}
			report.Errors = append(report.Errors, e)
			beego.Error("no draft created for user", uid, ":", e.Problem)
		} else {
			fbDraft, empty, err := inv.FastbillCreateDraft(false)
			if err == fastbill.ErrInvoiceAlreadyExported {
				beego.Info("draft for user", uid, "already exported")
				report.AlreadyExportedUids = append(report.AlreadyExportedUids, uid)
				report.SuccessUids = append(report.SuccessUids, uid)
				continue
			} else if err != nil {
				e := DraftsCreationError{
					UserId:  uid,
					Problem: err.Error(),
				}
				report.Errors = append(report.Errors, e)
				beego.Error("create draft for user", uid, ":", err)
				continue
			} else if !empty {
				report.SuccessUids = append(report.SuccessUids, uid)
			}
			if empty {
				report.EmptyUids = append(report.EmptyUids, uid)
				beego.Debug("draft is empty")
				continue
			}
			id := fbDraft.Id
			beego.Info("Draft created with ID", id)
			report.Ids = append(report.Ids, id)
		}
	}
	return
}
