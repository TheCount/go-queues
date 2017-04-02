package queue_test

import(
	"github.com/TheCount/go-queues/queue"
	"log"
	"sync"
)

// In this example, we build two integer queues
// and play around with them a little.

// IntQueue is our example structure to hold queues for integers.
type IntQueue struct {
	// Enqueue is identified as the enqueueing method by the
	// `queue:"enqueue"` structure tag.
	// This tag is sufficient to identify the enqueuing method.
	// We could have given the Enqueue method a different name,
	// but for clarity it makes sense to stick with Enqueue.
	// Enqueue takes exactly one argument,
	// the element x to be enqueued, and returns nothing.
	// Since this is supposed to be an integer queue,
	// we declare x to be of type int.
	Enqueue func( x int ) `queue:"enqueue"`

	// Analogous to Enqueue, Dequeue is identified as the dequeueing
	// method by the structure tag `queue:"dequeue"`.
	// Dequeue does not take any argument,
	// and returns two values: the dequeued element x and ok,
	// a boolean indicator whether the dequeue operation was successful.
	// A value of true for ok means success,
	// a value of false means that the queue was empty just as Dequeue
	// attempted to dequeue an element. In this case, x is the zero value,
	// which in the case of int is 0.
	// It would be an error to declare the return value x
	// to be of different type
	// than the argument x of the enqueueing method.
	Dequeue func()( x int, ok bool ) `queue:"dequeue"`
}

func Example() {
	// First, we make a queue for non-concurrent use.
	var ncq IntQueue
	if err := queue.Make( &ncq, queue.DefaultConfig().NonConcurrent() ); err != nil {
		log.Fatal( err )
	}
	// Let's enqueue a bunch of integers...
	for i := 1; i <= 10; i++ {
		ncq.Enqueue( i )
		log.Printf( "Enqueued %d\n", i )
	}
	// ...and dequeue them.
	for i := 1; i <= 10; i++ {
		x, _ := ncq.Dequeue()
		log.Printf( "Dequeued %d\n", x )
	}
	// The queue is now empty. Trying to dequeue again fails:
	_, ok := ncq.Dequeue()
	log.Printf( "Success dequeueing from empty queue: %v\n", ok )

	// Let's try concurrent queue access now.
	var ccq IntQueue
	if err := queue.Make( &ccq, nil /* passing nil selects a default implementation suitable for concurrent use */ ); err != nil {
		log.Fatal( err )
	}
	// We create two writing (enqueueing) and two reading goroutines.
	var wg sync.WaitGroup
	writer := func( id int ) {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			ccq.Enqueue( i )
			log.Printf( "Goroutine %d enqueued %d\n", id, i )
		}
	}
	reader := func( id int ) {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			x, ok := ccq.Dequeue()
			if ok {
				log.Printf( "Goroutine %d dequeued %d\n", id, x )
			} else {
				log.Printf( "Goroutine %d: queue was empty\n", id )
			}
		}
	}
	wg.Add( 4 )
	go writer( 1 )
	go writer( 2 )
	go reader( 3 )
	go reader( 4 )
	wg.Wait()
}
