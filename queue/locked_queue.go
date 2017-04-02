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
	dbFactory
	lq *lockedQueue
}

func ( lqf *lockedQueueFactory ) prepare() {
	lqf.lq = &lockedQueue{
		simpleQueue: *newSimpleQueue( lqf.capacityPerBuffer ),
		mx: sync.Mutex{},
	}
}

func ( lqf *lockedQueueFactory ) commit() {
	// empty
}

func ( lqf *lockedQueueFactory ) makeEnqueue( methodType reflect.Type ) reflect.Value {
	return makeEnqueue( lqf.lq, methodType )
}

func( lqf *lockedQueueFactory ) makeDequeue( methodType reflect.Type ) reflect.Value {
	return makeDequeue( lqf.lq, methodType )
}

func ( lqf *lockedQueueFactory ) reset() {
	lqf.lq = nil
}

func newLockedQueueFactory( initialCapacity int ) factory {
	return &lockedQueueFactory{
		dbFactory: *newDbFactory( initialCapacity ),
		lq: nil,
	}
}
