package memberstats

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/metrics/bin"
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
	b.UmsByName[k] = append(b.UmsByName[k], um)
}

func newBins(from, to month.Month) (bs []*Bin) {
	last, ok := bin.Map(from, to, to)
	if !ok {
		panic(fmt.Sprintf("%v/%v/%v", from, to, to))
	}

	return make([]*Bin, last+1)
}

func (s *Stats) Bins() (bs []*Bin) {
	bs = newBins(s.from, s.to)

	for _, iv := range s.invs {
		i, ok := bin.Map(s.from, s.to, iv.GetMonth())
		if !ok {
			continue
		}
		for _, ium := range iv.InvUserMemberships {
			if bs[i] == nil {
				bs[i] = NewBin(iv.GetMonth())
			}
			bs[i].Add(ium.UserMembership)
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
