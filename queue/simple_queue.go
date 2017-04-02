/*
Copyright (c) 2017 Alexander Klauer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
	dbFactory
	sq *simpleQueue
}

func newSimpleQueue( capacityPerBuffer int ) *simpleQueue {
	return &simpleQueue{
		buf1: make( []interface{}, capacityPerBuffer ),
		buf2: make( []interface{}, capacityPerBuffer ),
		start: 0,
		end: 0,
	}
}

func ( sqf *simpleQueueFactory ) prepare() {
	sqf.sq = newSimpleQueue( sqf.capacityPerBuffer )
}

func ( sqf *simpleQueueFactory ) commit() {
	// empty
}

func ( sqf *simpleQueueFactory ) makeEnqueue( methodType reflect.Type ) reflect.Value {
	return makeEnqueue( sqf.sq, methodType )
}

func( sqf *simpleQueueFactory ) makeDequeue( methodType reflect.Type ) reflect.Value {
	return makeDequeue( sqf.sq, methodType )
}

func ( sqf *simpleQueueFactory ) reset() {
	sqf.sq = nil
}

func newSimpleQueueFactory( initialCapacity int ) factory {
	return &simpleQueueFactory{
		dbFactory: *newDbFactory( initialCapacity ),
		sq: nil,
	}
}
