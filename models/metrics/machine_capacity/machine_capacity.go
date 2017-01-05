package machine_capacity

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
	"github.com/FabLabBerlin/localmachines/models/metrics/filter"
	"github.com/FabLabBerlin/localmachines/models/purchases"
	"time"
)

type MachineCapacity struct {
	m    *machine.Machine
	from month.Month
	to   month.Month
	invs []*invutil.Invoice
}

func New(
	m *machine.Machine,
	from month.Month,
	to month.Month,
	invs []*invutil.Invoice,
) *MachineCapacity {

	return &MachineCapacity{
		m:    m,
		from: from,
		to:   to,
		invs: filter.Invoices(invs, from, to),
	}
}

// Percentage is a number between 0 and 1
type Percentage float64

func (mc MachineCapacity) Capacity() time.Duration {
	to := day.New(mc.to.Year(), mc.to.Month(), 1).AddDate(0, 1, -1)
	dur := to.Sub(mc.Opening())
	if dur.Seconds() < 0 {
		return 0
	}
	return dur
}

func (mc MachineCapacity) CapacityCached() (c time.Duration, err error) {
	key := fmt.Sprintf("Capacity(%v)-%v-%v", mc.m.Id, mc.from, mc.to)

	err = redis.Cached(key, 3600, &c, func() (interface{}, error) {
		return mc.Capacity(), nil
	})

	return
}

func (mc MachineCapacity) ContainsTime(t time.Time) bool {
	return !mc.from.AfterTime(t) && !mc.to.BeforeTime(t)
}

func (mc MachineCapacity) Opening() (opening day.Day) {
	for _, inv := range mc.invs {
		for _, p := range inv.Purchases {
			if p.MachineId != mc.m.Id {
				continue
			}

			if !mc.ContainsTime(p.TimeStart) {
				continue
			}

			if opening.IsZero() || opening.AfterTime(p.TimeStart) {
				opening = day.NewTime(p.TimeStart)
			}
		}
	}

	if opening.IsZero() {
		opening = day.Now()
	}

	return
}

func (mc MachineCapacity) Usage() (usage time.Duration) {
	for _, inv := range mc.invs {
		for _, p := range inv.Purchases {
			if p.MachineId != mc.m.Id {
				continue
			}
			if p.Type != purchases.TYPE_ACTIVATION {
				continue
			}
			if !mc.ContainsTime(p.TimeStart) {
				continue
			}

			dur := time.Duration(p.Seconds()) * time.Second
			usage += dur
		}
	}

	return
}

func (mc MachineCapacity) UsageCached() (u time.Duration, err error) {
	key := fmt.Sprintf("Usage(%v)-%v-%v", mc.m.Id, mc.from, mc.to)

	err = redis.Cached(key, 3600, &u, func() (interface{}, error) {
		return mc.Usage(), nil
	})

	return
}

func (mc MachineCapacity) Utilization() Percentage {
	usage := mc.Usage().Seconds()

	if usage == 0 {
		return 0
	}

	return Percentage(usage / mc.Capacity().Seconds())
}

func (mc MachineCapacity) UtilizationCached() (u Percentage, err error) {
	key := fmt.Sprintf("Utilization(%v)-%v-%v", mc.m.Id, mc.from, mc.to)

	err = redis.Cached(key, 3600, &u, func() (interface{}, error) {
		return mc.Utilization(), nil
	})

	return
}
