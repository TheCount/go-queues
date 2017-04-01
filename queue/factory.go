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
