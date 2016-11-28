package retention

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"github.com/FabLabBerlin/localmachines/models/users"
	"time"
)

type Retention struct {
	locationId int64
	stepDays   int
	from       day.Day
	to         day.Day
	invs       []*invutil.Invoice
	us         []*users.User
}

func New(
	locId int64,
	stepDays int,
	from day.Day,
	to day.Day,
	invs []*invutil.Invoice,
	us []*users.User,
) *Retention {

	return &Retention{
		locationId: locId,
		stepDays:   stepDays,
		from:       from,
		to:         to,
		invs:       invs,
		us:         us,
	}
}

type Percent float64

// Row is typically a day or a week
type Row struct {
	From          day.Day
	StepDays      int
	newUsers      []int64
	Returns       []Percent
	returnedUsers [][]int64
	Users         int
}

func NewRow(from day.Day, stepDays, numberOfValues int) *Row {
	r := &Row{
		From:          from,
		StepDays:      stepDays,
		newUsers:      make([]int64, 0, 20),
		Returns:       make([]Percent, numberOfValues),
		returnedUsers: make([][]int64, numberOfValues),
	}

	for i := range r.returnedUsers {
		r.returnedUsers[i] = make([]int64, 0, 20)
	}

	return r
}

func (row *Row) AddNewUser(userId int64) {
	fmt.Printf("row[from=%v]#AddNewUser(%v)\n", row.From, userId)
	row.newUsers = append(row.newUsers, userId)
}

func (row *Row) AddReturnedUser(t time.Time, userId int64) {
	i := dayToIndex(row.From, day.NewTime(t), row.StepDays)

	row.returnedUsers[i] = append(row.returnedUsers[i], userId)
}

func (row *Row) Calculate() {
	var h map[int64]struct{}
	row.newUsers, h = uniq(row.newUsers)
	row.Users = len(row.newUsers)

	for i, returned := range row.returnedUsers {
		row.returnedUsers[i], _ = uniq(returned)

		n := 0

		for _, id := range row.returnedUsers[i] {
			if _, ok := h[id]; ok {
				n++
			}
		}

		if n > 0 {
			row.Returns[i] = Percent(float64(n) / float64(row.Users))
		}
	}
}

func (row Row) NewUsers() []int64 {
	return row.newUsers
}

func (row Row) ReturnedUsers() [][]int64 {
	return row.returnedUsers
}

func (row Row) To() day.Day {
	return row.From.AddDate(0, 0, row.StepDays-1)
}

func uniq(ids []int64) (u []int64, h map[int64]struct{}) {
	h = make(map[int64]struct{})

	for _, id := range ids {
		h[id] = struct{}{}
	}

	u = make([]int64, 0, len(h))

	for id := range h {
		u = append(u, id)
	}

	return
}

// Calculate a retention triangle like in Mixpanel or Google Analytics.
func (r Retention) Calculate() (triangle []*Row) {
	triangle = make([]*Row, r.NumberOfRows()-1)

	i := 0
	for d := r.from; d.BeforeOrEqual(r.to); d = d.AddDate(0, 0, int(r.stepDays)) {
		triangle[i] = NewRow(d, r.stepDays, r.NumberOfRows()-i-1)

		if i++; i >= len(triangle) {
			break
		}
	}

	for _, inv := range r.invs {
		if inv.Canceled {
			continue
		}

		for _, p := range inv.Purchases {
			i := r.RowFor(p)
			if i >= len(triangle) {
				continue
			}
			fmt.Printf("i=%v p.TimeStart=%v\n", i, p.TimeStart)
			for j := 0; j <= i; j++ {
				fmt.Printf("[j=%v] \n", j)
				triangle[j].AddReturnedUser(p.TimeStart, p.UserId)
			}
		}
	}

	fmt.Printf("len(r.us)=%v\n", len(r.us))
	for _, u := range r.us {
		for _, row := range triangle {
			d := day.NewTime(u.Created)
			if row.From.BeforeOrEqual(d) && d.BeforeOrEqual(row.To()) {
				row.AddNewUser(u.Id)
				row.Calculate()
				break
			}
		}

	}

	return
}

func (r Retention) NumberOfRows() int {
	return r.RowForDay(r.to) + 1
}

func (r Retention) RowFor(p *purchases.Purchase) int {
	return r.RowForDay(day.NewTime(p.TimeStart))
}

func (r Retention) RowForDay(d day.Day) int {
	return dayToIndex(r.from, d, r.stepDays)
}

func dayToIndex(d0, d day.Day, stepDays int) int {
	f := d.Sub(d0).Hours() / 24.0 / float64(stepDays)

	return int(f)
}
