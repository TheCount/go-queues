package queue

// T is a placeholder for an actual queue element type.
type T interface{}

// GenericQueue is a template for a queue structure.
// You can copy and paste this structure,
// give it a name of your choosing,
// and then replace T with your element type.
// A pointer to an instance of your new structure can then be passed
// to the Make() function.
type GenericQueue struct {
	// Enqueue enqueues element x into the queue.
	// The tag is used by Make()
	// to identify this as the enqueueing function.
	// If you like, you can give this function a different name.
	Enqueue func( x T ) `queue:"enqueue"`

	// Dequeue attempts to dequeue an element from the queue.
	// If successful, the dequeued element is returned as x
	// and ok is true.
	// If unsuccessful, i. e.,
	// the queue was empty
	// at the time the dequeueing operation was attempted,
	// the value of x is the zero value of T and ok is false.
	// The tag is used by Make()
	// to identify this as the dequeueing function.
	// If you like, you can give this function a different name.
	Dequeue func()( x T, ok bool ) `queue:"dequeue"`
}

// interfaceQueue is the minimal generic queue interface used internally.
type interfaceQueue interface {
	enqueue( x interface{} )
	dequeue() ( x interface{}, ok bool )
}
