package queue

import(
	"reflect"
	"sync"
)

// lockedQueue uses a mutex to make simpleQueue totally thread-safe.
type lockedQueue struct {
	simpleQueue
	mx sync.Mutex
}

func ( q *lockedQueue ) enqueue( x interface{} ) {
	q.mx.Lock()
	defer q.mx.Unlock()
	q.simpleQueue.enqueue( x )
}

func ( q *lockedQueue ) dequeue() ( interface{}, bool ) {
	q.mx.Lock()
	defer q.mx.Unlock()
	return q.simpleQueue.dequeue()
}

// lockedQueueFactory implements factory for lockedQueue
type lockedQueueFactory struct {
	dbFactory
	lq *lockedQueue
}

func ( lqf *lockedQueueFactory ) prepare() {
	lqf.lq = &lockedQueue{
		simpleQueue: *newSimpleQueue( lqf.capacityPerBuffer ),
		mx: sync.Mutex{},
	}
}

func ( lqf *lockedQueueFactory ) commit() {
	// empty
}

func ( lqf *lockedQueueFactory ) makeEnqueue( methodType reflect.Type ) reflect.Value {
	return makeEnqueue( lqf.lq, methodType )
}

func( lqf *lockedQueueFactory ) makeDequeue( methodType reflect.Type ) reflect.Value {
	return makeDequeue( lqf.lq, methodType )
}

func ( lqf *lockedQueueFactory ) reset() {
	lqf.lq = nil
}

func newLockedQueueFactory( initialCapacity int ) factory {
	return &lockedQueueFactory{
		dbFactory: *newDbFactory( initialCapacity ),
		lq: nil,
	}
}
