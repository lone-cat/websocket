package sem

import "sync"

func isSignalChannelClosed(ch chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

type TwoStage struct {
	mu           sync.RWMutex
	stoppingChan chan struct{}
	stoppedChan  chan struct{}
}

func (tss *TwoStage) init() {
	var needInit bool

	tss.mu.RLock()
	if tss.stoppingChan == nil || tss.stoppedChan == nil {
		needInit = true
	}
	tss.mu.RUnlock()

	if needInit {
		tss.mu.Lock()
		tss.stoppingChan = make(chan struct{})
		tss.stoppedChan = make(chan struct{})
		tss.mu.Unlock()
		close(tss.stoppingChan)
		close(tss.stoppedChan)
	}
}

func (tss *TwoStage) IsStopping() bool {
	tss.init()

	return isSignalChannelClosed(tss.stoppingChan)
}

func (tss *TwoStage) IsStopped() bool {
	tss.init()

	return isSignalChannelClosed(tss.stoppedChan)
}

func (tss *TwoStage) Start() bool {
	tss.init()

	if !tss.IsStopped() {
		return false
	}

	tss.mu.Lock()
	defer tss.mu.Unlock()

	tss.stoppingChan = make(chan struct{})
	tss.stoppedChan = make(chan struct{})

	return true
}

func (tss *TwoStage) StartStopping() {
	tss.init()

	if !tss.IsStopping() {
		tss.mu.Lock()
		defer tss.mu.Unlock()
		close(tss.stoppingChan)
	}
}

func (tss *TwoStage) FinishStopping() {
	tss.init()

	tss.StartStopping()

	if !tss.IsStopped() {
		tss.mu.Lock()
		defer tss.mu.Unlock()
		close(tss.stoppedChan)
	}
}

func (tss *TwoStage) WaitTillStopped() {
	tss.init()

	<-tss.GetStoppedChannel()
}

func (tss *TwoStage) GetStoppingChannel() chan struct{} {
	tss.init()
	tss.mu.RLock()
	defer tss.mu.RUnlock()
	stoppingChan := tss.stoppingChan
	return stoppingChan
}

func (tss *TwoStage) GetStoppedChannel() chan struct{} {
	tss.init()
	tss.mu.RLock()
	defer tss.mu.RUnlock()
	stopChan := tss.stoppedChan
	return stopChan
}
