package listener

import (
	"time"
)

type Factory struct {
}

func (f Factory) CreateDebouncer(debounce time.Duration, semaphore StopSemaphoreI, logger LoggerI) *Debouncer {
	return NewDebouncer(debounce, semaphore, logger)
}

func (f Factory) CreateListener(local bool, port uint16, semaphore StopSemaphoreI, logger LoggerI) *Listener {
	return NewListener(local, port, semaphore, logger)
}

func (f Factory) CreateLimiter(connectionsLimit uint8, stopSemaphore StopSemaphoreI, logger LoggerI) *Limiter {
	return NewLimiter(connectionsLimit, stopSemaphore, logger)
}

func (f Factory) CreateConnectionProvider(
	local bool,
	port uint16,
	listenerSem StopSemaphoreI,
	listenerLogger LoggerI,
	debounce time.Duration,
	debouncerSem StopSemaphoreI,
	debouncerLogger LoggerI,
	connectionsLimit uint8,
	limiterSem StopSemaphoreI,
	limiterLogger LoggerI,
) *ConnectionProvider {
	listener := f.CreateListener(local, port, listenerSem, listenerLogger)
	debouncer := f.CreateDebouncer(debounce, debouncerSem, debouncerLogger)
	limiter := f.CreateLimiter(connectionsLimit, limiterSem, limiterLogger)
	return NewConnectionProvider(listener, debouncer, limiter)
}
