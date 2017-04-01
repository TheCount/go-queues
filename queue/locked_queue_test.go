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
	var wg sync.WaitGroup
	writer := func() {
		defer wg.Done()
		for i := 0; i != 10000; i++ {
			enqueue( i )
		}
	}
	reader := func() {
		defer wg.Done()
		previous := -1
		minimum := 0
		for i := 0; i != 100000; i++ {
			x, ok := dequeue()
			if ( previous == 100 ) && ( minimum == 100 ) && ok {
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
