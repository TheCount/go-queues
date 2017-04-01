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
	lq *lockedQueue
	capacityPerBuffer int
}

func ( lqf *lockedQueueFactory ) prepare() {
	lqf.lq = &lockedQueue{
		simpleQueue: simpleQueue{
			buf1: make( []interface{}, lqf.capacityPerBuffer ),
			buf2: make( []interface{}, lqf.capacityPerBuffer ),
			start: 0,
			end: 0,
		},
		mx: sync.Mutex{},
	}
}

func ( lqf *lockedQueueFactory ) commit() {
	// empty
}

func ( lqf *lockedQueueFactory ) makeEnqueue( methodType reflect.Type ) reflect.Value {
	q := lqf.lq
	return reflect.MakeFunc( methodType, func( args []reflect.Value ) []reflect.Value {
		q.enqueue( args[0].Interface() )
		return []reflect.Value{}
	} )
}

func( lqf *lockedQueueFactory ) makeDequeue( methodType reflect.Type ) reflect.Value {
	q := lqf.lq
	return reflect.MakeFunc( methodType, func( args []reflect.Value ) []reflect.Value {
		x, ok := q.dequeue()
		if ok {
			return []reflect.Value{
				reflect.ValueOf( x ),
				reflect.ValueOf( ok ),
			}
		} else {
			return []reflect.Value{
				reflect.Zero( methodType.Out( 0 ) ),
				reflect.ValueOf( ok ),
			}
		}
	} )
}

func ( lqf *lockedQueueFactory ) reset() {
	lqf.lq = nil
}

func newLockedQueueFactory( initialCapacity int ) factory {
	capacityPerBuffer := initialCapacity / 2
	if capacityPerBuffer < 1 {
		capacityPerBuffer = 1
	}

	return &lockedQueueFactory{
		lq: nil,
		capacityPerBuffer: capacityPerBuffer,
	}
}
