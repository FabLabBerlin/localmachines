package machine_capacity

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/lib/day"
	"github.com/FabLabBerlin/localmachines/lib/month"
	"github.com/FabLabBerlin/localmachines/lib/redis"
	"github.com/FabLabBerlin/localmachines/models/invoices/invutil"
	"github.com/FabLabBerlin/localmachines/models/machine"
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
		invs: invs,
	}
}

// Percentage is a number between 0 and 1
type Percentage float64

func (mc MachineCapacity) Capacity() time.Duration {
	return day.Now().Sub(mc.Opening())
}

func (mc MachineCapacity) CapacityCached() (c time.Duration) {
	key := fmt.Sprintf("Capacity(%v)", mc.m.Id)

	redis.Cached(key, 3600, &c, func() interface{} {
		return mc.Capacity()
	})

	return
}

func (mc MachineCapacity) Opening() (opening day.Day) {
	for _, inv := range mc.invs {
		for _, p := range inv.Purchases {
			if p.MachineId != mc.m.Id {
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
		if inv.Canceled {
			continue
		}

		for _, p := range inv.Purchases {
			if p.MachineId != mc.m.Id {
				continue
			}
			if p.Type != purchases.TYPE_ACTIVATION {
				continue
			}

			dur := time.Duration(p.Seconds()) * time.Second
			usage += dur
		}
	}

	return
}

func (mc MachineCapacity) UsageCached() (u time.Duration) {
	key := fmt.Sprintf("Usage(%v)", mc.m.Id)

	redis.Cached(key, 3600, &u, func() interface{} {
		return mc.Usage()
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

func (mc MachineCapacity) UtilizationCached() (u Percentage) {
	key := fmt.Sprintf("Utilization(%v)", mc.m.Id)

	redis.Cached(key, 3600, &u, func() interface{} {
		return mc.Utilization()
	})

	return
}
