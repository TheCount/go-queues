package queue

import(
	"reflect"
)

// simpleQueue keeps the data for a simple, non-concurrent queue.
type simpleQueue struct {
	buf1 []interface{}
	buf2 []interface{}
	start, end int
}

func ( q *simpleQueue ) enqueue( x interface{} ) {
	if q.end >= len( q.buf1 ) {
		if q.end >= len( q.buf1 ) + len( q.buf2 ) {
			q.buf2 = append( q.buf2, x )
		} else {
			q.buf2[q.end - len( q.buf1 )] = x
		}
	} else {
		q.buf1[q.end] = x
	}
	q.end++
}

func ( q *simpleQueue ) dequeue() ( x interface{}, ok bool ) {
	if q.start == q.end {
		ok = false
		return
	}
	x = q.buf1[q.start]
	q.buf1[q.start] = nil
	ok = true
	q.start++
	if q.start == len( q.buf1 ) {
		q.start -= len( q.buf1 )
		q.end -= len( q.buf1 )
		q.buf1, q.buf2 = q.buf2, q.buf1
	}

	return
}

// simpleQueueFactory implements factory for simpleQueue
type simpleQueueFactory struct {
	sq *simpleQueue
	capacityPerBuffer int
}

func ( sqf *simpleQueueFactory ) prepare() {
	sqf.sq = &simpleQueue{
		buf1: make( []interface{}, sqf.capacityPerBuffer ),
		buf2: make( []interface{}, sqf.capacityPerBuffer ),
		start: 0,
		end: 0,
	}
}

func ( sqf *simpleQueueFactory ) commit() {
	// empty
}

func ( sqf *simpleQueueFactory ) makeEnqueue( methodType reflect.Type ) reflect.Value {
	q := sqf.sq
	return reflect.MakeFunc( methodType, func( args []reflect.Value ) []reflect.Value {
		q.enqueue( args[0].Interface() )
		return []reflect.Value{}
	} )
}

func( sqf *simpleQueueFactory ) makeDequeue( methodType reflect.Type ) reflect.Value {
	q := sqf.sq
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

func ( sqf *simpleQueueFactory ) reset() {
	sqf.sq = nil
}

func newSimpleQueueFactory( initialCapacity int ) factory {
	capacityPerBuffer := initialCapacity / 2
	if capacityPerBuffer < 1 {
		capacityPerBuffer = 1
	}

	return &simpleQueueFactory{
		sq: nil,
		capacityPerBuffer: capacityPerBuffer,
	}
}
