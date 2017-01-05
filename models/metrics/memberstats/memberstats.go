package memberstats

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/metrics/filter"
	"github.com/FabLabBerlin/localmachines/models/user_memberships"
	"time"
)

type Stats struct {
	from month.Month
	to   month.Month
	invs []*invutil.Invoice
}

func New(
	from month.Month,
	to month.Month,
	invs []*invutil.Invoice,
) *Stats {

	return &Stats{
		from: from,
		to:   to,
		invs: filter.Invoices(invs, from, to),
	}
}

type Bin struct {
	Month     month.Month
	UmsByName map[string][]*user_memberships.UserMembership
}

func NewBin(m month.Month) *Bin {
	return &Bin{
		Month:     m,
		UmsByName: make(map[string][]*user_memberships.UserMembership),
	}
}

func (b *Bin) Add(um *user_memberships.UserMembership) {
	k := um.Membership.Title
	if _, ok := b.UmsByName[k]; !ok {
		b.UmsByName[k] = make([]*user_memberships.UserMembership, 0, 40)
	}
	fmt.Printf("k=%v\n", k)
	b.UmsByName[k] = append(b.UmsByName[k], um)
}

func newBins(from, to month.Month) (bins []*Bin) {
	n := 0
	for t := from; !t.After(to); t = t.Add(1) {
		n++
	}
	bins = make([]*Bin, n)

	return
}

func MapBin(from, to, m month.Month) (i int, ok bool) {
	for t := from; t.BeforeOrEqual(to); t = t.Add(1) {
		if t.Equal(m) {
			return i, true
		}
		i++
	}
	return -1, false
}

func (s *Stats) Bins() (bins []*Bin) {
	bins = newBins(s.from, s.to)

	for _, iv := range s.invs {
		i, ok := MapBin(s.from, s.to, iv.GetMonth())
		if !ok {
			fmt.Printf("oh no, couldn't map where s.from=%v and s.to=%v: %v\n", s.from, s.to, iv.GetMonth())
			continue
		}
		for _, ium := range iv.InvUserMemberships {
			if bins[i] == nil {
				bins[i] = NewBin(iv.GetMonth())
			}
			bins[i].Add(ium.UserMembership)
		}
	}

	return
}

func (s *Stats) BinsCached() (c time.Duration, err error) {
	key := fmt.Sprintf("Memberstats-%v-%v", s.from, s.to)

	err = redis.Cached(key, 3600, &c, func() (interface{}, error) {
		return s.Bins(), nil
	})

	return
}
