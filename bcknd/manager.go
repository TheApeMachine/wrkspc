package bcknd

import (
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/nemo"
	"github.com/theapemachine/wrkspc/sockpuppet"
	"github.com/theapemachine/wrkspc/spdg"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Manager orchestrates the functionality of a service and brings the results to Egress.
*/
type Manager struct {
	egress   *Egress
	Hub      *sockpuppet.Hub
	stores   []Store
	disposer *twoface.Disposer
}

/*
NewManager constructs a Manager and returns a reference to it.
*/
func NewManager(egress *Egress, hub *sockpuppet.Hub, disposer *twoface.Disposer) *Manager {
	errnie.Traces()

	return &Manager{
		egress:   egress,
		Hub:      hub,
		disposer: disposer,
		stores: []Store{
			NewStore(nemo.NewClient(
				viper.GetString("wrkspc.nemo.access-key-id"),
				viper.GetString("wrkspc.nemo.access-key-secret"),
				viper.GetString("wrkspc.nemo.region"),
				viper.GetString("wrkspc.nemo.bucket"),
			)),
		},
	}
}

/*
Execute the job the Manager is tasked with by employing workers.
*/
func (manager *Manager) Execute(question *spdg.Datagram) chan *spdg.Datagram {
	errnie.Traces()

	out := make(chan *spdg.Datagram)
	agg := make(chan *spdg.Datagram)

	for _, store := range manager.stores {
		manager.search(store, question, agg)
	}

	return out
}

/*
search a Store.
*/
func (manager *Manager) search(store Store, question *spdg.Datagram, agg chan *spdg.Datagram) {
	errnie.Traces()

	go func() {
		for {
			select {
			case dg := <-store.Peek(question):
				agg <- dg
			case <-manager.disposer.Done():
				return
			}
		}
	}()
}
