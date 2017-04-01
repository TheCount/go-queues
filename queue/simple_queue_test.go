package queue

import(
	"reflect"
	"testing"
)

func TestCapacity( t *testing.T ) {
	for i := -3; i < 4; i++ {
		f := newSimpleQueueFactory( i ).( *simpleQueueFactory )
		if f.capacityPerBuffer < 1 {
			t.Errorf( "Too low buffer capacity: %d (requested: %d)", f.capacityPerBuffer, i )
		}
	}
}

func TestSimpleQueue ( t *testing.T ) {
	// Build queue
	f := newSimpleQueueFactory( 0 )
	f.prepare()
	var enqueue func( int )
	var dequeue func() ( int, bool )
	enqueue = f.makeEnqueue( reflect.TypeOf( enqueue ) ).Interface().( func( int ) )
	dequeue = f.makeDequeue( reflect.TypeOf( dequeue ) ).Interface().( func() ( int, bool ) )
	f.commit()
	f.reset()
	// Empty dequeue check
	x, ok := dequeue()
	if ok {
		t.Error( "Dequeue succeeds on empty queue" )
	}
	if x != 0 {
		t.Errorf( "Failed dequeue does not return zero value: %d", x )
	}
	// Single enqueue, double dequeue check
	enqueue( 42 )
	x, ok = dequeue()
	if !ok {
		t.Error( "Dequeue on non-empty queue failed" )
	}
	if x != 42 {
		t.Errorf( "Dequeue returned wrong value: %d instead of 42", x )
	}
	x, ok = dequeue()
	if ok {
		t.Error( "Dequeue succeeds on now-empty queue" )
	}
	if x != 0 {
		t.Errorf( "Failed dequeue does not return zero value: %d", x )
	}
	// queue order check
	for i := 0; i < 100; i++ {
		enqueue( i )
	}
	for i := 0; i < 100; i++ {
		x, ok = dequeue()
		if !ok {
			t.Error( "Dequeue fails on non-empty queue" )
		}
		if x != i {
			t.Errorf( "Dequeue returned wrong value: %d instead of %d", x, i )
		}
	}
}