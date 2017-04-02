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
