package queue

import(
	"reflect"
)

// factory is an internal interface
// helping with building a specific queue type.
type factory interface {
	// prepare prepares the factory for creation of a new queue
	prepare()

	// commit commits preparations and makes the queue usable
	commit()

	// makeEnqueue creates the enqueueing method
	makeEnqueue( methodType reflect.Type ) reflect.Value

	// makeDequeue creates the dequeueing method
	makeDequeue( methodType reflect.Type ) reflect.Value

	// reset resets preparations without committing them.
	// Calling reset before prepare() or after commit() has no effect.
	reset()
}

// dbFactory is the basic building block for the factories of double-buffered queues.
type dbFactory struct {
	// capacityPerBuffer is the initial capacity of each of the queue buffers.
	// It must be at least 1.
	capacityPerBuffer int
}

// newDbFactory creates a new dbFactory for double-buffered queues.
// The argument initialCapacity is the total initial capacity of the queue.
// Values too small will be corrected.
func newDbFactory( initialCapacity int ) *dbFactory {
	capacityPerBuffer := initialCapacity / 2
	if capacityPerBuffer < 1 {
		capacityPerBuffer = 1
	}

	return &dbFactory{
		capacityPerBuffer: capacityPerBuffer,
	}
}

// makeEnqueue creates the function interfacing the typed enqueue function
// with the generic implementation of the queue.
func makeEnqueue( q interfaceQueue, methodType reflect.Type ) reflect.Value {
	return reflect.MakeFunc( methodType, func( args []reflect.Value ) []reflect.Value {
		q.enqueue( args[0].Interface() )
		return []reflect.Value{}
	} )
}

// makeDequeue creates the function interfacing the typed dequeue function
// with the generic implementation of the queue.
func makeDequeue( q interfaceQueue, methodType reflect.Type ) reflect.Value {
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
