package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"time"
)

type Invoice struct {
	Id          int64  `orm:"auto";"pk"`
	Activations string `orm:type(text)`
	XlsFile     string `orm:size(255)`
}

func (this *Invoice) TableName() string {
	return "invoices"
}

func init() {
	orm.RegisterModel(new(Invoice))
}

func CreateInvoice(startTime time.Time,
	endTime time.Time) (*Invoice, error) {

	beego.Info("Creating invoice...")

	act := Activation{}
	usr := User{}
	mch := Machine{}

	o := orm.NewOrm()

	// Get unique users
	users := []User{}
	query := fmt.Sprintf("SELECT DISTINCT u.* FROM %s u JOIN %s a ON a.user_id=u.id "+
		"WHERE a.time_start>? AND a.time_end<? AND a.invoiced=false AND a.active=false ",
		usr.TableName(),
		act.TableName())

	num, err := o.Raw(query,
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02")).QueryRows(&users)

	if err != nil {
		msg := fmt.Sprintf("Failed to get unique users: %v", err)
		return nil, errors.New(msg)
	}

	beego.Trace("Num uniq users:", num)

	// Create a xlsx file if there
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	if len(users) > 0 {

		file = xlsx.NewFile()
		sheet = file.AddSheet("Sheet1")
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "I am a cell!"

	} else {
		return nil, errors.New("There are no invoiceable activations")
	}

	// Loop through users
	for _, user := range users {

		// 1. Get unique user machines within the range of activations
		machines := []Machine{}
		query = fmt.Sprintf("SELECT DISTINCT m.* FROM %s m JOIN %s a ON a.machine_id=m.id "+
			"WHERE a.time_start>? AND a.time_end<? AND a.invoiced=false "+
			"AND a.active=false AND a.user_id=?",
			mch.TableName(),
			act.TableName())

		num, err = o.Raw(query,
			startTime.Format("2006-01-02"),
			endTime.Format("2006-01-02"),
			user.Id).QueryRows(&machines)

		if err != nil {
			msg := fmt.Sprintf("Failed to get user machines: %v", err)
			return nil, errors.New(msg)
		}

		beego.Trace("Num uniq machines:", num)

		if len(machines) > 0 {
			beego.Trace("==========")
			beego.Trace(user)
			beego.Trace(machines)

			// 2. Loop through user machines and get activations per user machine
			for _, machine := range machines {

				var sum int64

				// Get activations for user and machine, get total time and other
				query = fmt.Sprintf("SELECT SUM(time_total) FROM %s "+
					"WHERE time_start>? AND time_end<? AND invoiced=false "+
					"AND active=false AND user_id=? AND machine_id=?",
					act.TableName())

				err = o.Raw(query,
					startTime.Format("2006-01-02"),
					endTime.Format("2006-01-02"),
					user.Id,
					machine.Id).QueryRow(&sum)

				if err != nil {
					msg := fmt.Sprintf("Failed to get activation sum: %v", err)
					return nil, errors.New(msg)
				}

				var final float32
				var finalUnits string
				if machine.PriceUnit == "minute" {
					final = float32(sum) / 60.0
					finalUnits = "minutes"
				} else if machine.PriceUnit == "hour" {
					final = float32(sum) / 60.0 / 60.0
					finalUnits = "hours"
				}

				var price float32
				price = machine.Price * final

				beego.Trace("-----")
				beego.Trace(machine.Name)
				beego.Trace(final, finalUnits)
				beego.Trace("price:", price)
			}

		}
	}

	filename := "files/MyXLSXFile.xlsx"
	err = file.Save(filename)
	if err != nil {
		fmt.Printf(err.Error())
	}

	inv := Invoice{}
	inv.XlsFile = filename

	return &inv, nil
}
