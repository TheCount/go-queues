package queue

import(
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestLockedQueue( t *testing.T ) {
	f := newLockedQueueFactory( 0 )
	f.prepare()
	var enqueue func( int )
	var dequeue func() ( int, bool )
	enqueue = f.makeEnqueue( reflect.TypeOf( enqueue ) ).Interface().( func( int ) )
	dequeue = f.makeDequeue( reflect.TypeOf( dequeue ) ).Interface().( func() ( int, bool ) )
	f.commit()
	f.reset()
	// Parallel queue check
	const iterations = 10000
	var wg sync.WaitGroup
	writer := func() {
		defer wg.Done()
		for i := 1; i <= iterations; i++ {
			enqueue( i )
		}
	}
	reader := func() {
		defer wg.Done()
		previous := 0
		minimum := 0
		for i := 1; i <= 10 * iterations; i++ {
			x, ok := dequeue()
			if ( previous == iterations ) && ( minimum == iterations ) && ok {
				t.Errorf( "Spurious successful dequeue: %d", x )
			}
			if !ok {
				runtime.Gosched()
				continue
			}
			if x < minimum {
				t.Errorf( "Out of order dequeue: %d", x )
			}
			if x <= previous {
				minimum = x
			}
			previous = x
		}
	}
	wg.Add( 4 )
	go writer()
	go writer()
	go reader()
	go reader()
	wg.Wait()
}
