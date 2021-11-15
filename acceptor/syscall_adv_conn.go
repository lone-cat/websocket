package acceptor

import (
	"errors"
	"sync"

	"github.com/mailru/easygo/netpoll"
)

type SyscallAdvancedConn struct {
	AdvancedNetConnI
	desc        *netpoll.Desc
	poller      netpoll.Poller
	mu          sync.RWMutex
	inited      bool
	closingFunc func()
}

func ConvertAdvancedConnToSyscall(conn AdvancedNetConnI, desc *netpoll.Desc, poller netpoll.Poller) *SyscallAdvancedConn {
	return &SyscallAdvancedConn{
		AdvancedNetConnI: conn,
		desc:             desc,
		poller:           poller,
	}
}

func (sac *SyscallAdvancedConn) Init(actionSuccess func()) (err error) {
	sac.mu.Lock()
	defer sac.mu.Unlock()
	realAction := func(event netpoll.Event) {
		if event&netpoll.EventHup == 0 && event&netpoll.EventReadHup == 0 && event&netpoll.EventErr == 0 {
			actionSuccess()
		} else {
			// handle error here
		}
	}

	err = sac.poller.Start(sac.desc, realAction)
	if err != nil {
		return
	}

	sac.closingFunc = func() {
		sac.poller.Stop(sac.desc)
		sac.desc.Close()
	}

	err = sac.poller.Stop(sac.desc)
	if err != nil {
		return
	}

	sac.inited = true

	return
}

func (sac *SyscallAdvancedConn) Resume() error {
	sac.mu.RLock()
	defer sac.mu.RUnlock()
	if !sac.inited {
		return errors.New(`not inited`)
	}
	return sac.poller.Resume(sac.desc)
}

func (sac *SyscallAdvancedConn) Close() error {
	sac.mu.Lock()
	defer sac.mu.Unlock()
	if sac.closingFunc != nil {
		sac.closingFunc()
		sac.closingFunc = nil
	}
	return sac.AdvancedNetConnI.Close()
}
