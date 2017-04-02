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
	"errors"
	"fmt"
	"reflect"
)

// Make creates a new queue.
// The argument qptr must be a pointer to an instance of
// a structure satisfying the constraints documented in GenericQueue.
// The parameter config can be used to specify the characteristics
// of the queue.
// A nil argument is permissible.
// In this case, the default configuration is used.
// On success, nil is returned.
// On error, an appropriate error is returned.
func Make( qptr interface{}, config *Config ) error {
	// Tag IDs
	const(
		queue = "queue"
		enqueue = "enqueue"
		dequeue = "dequeue"
	)
	// Get config
	if config == nil {
		config = DefaultConfig()
	}
	if !config.IsValid() {
		return errors.New( "Invalid queue configuration" )
	}
	factory := config.factory()
	if factory == nil {
		return errors.New( "This queue configuration has not been implemented yet" )
	}
	factory.prepare()
	defer factory.reset()
	// Extract function pointers
	qptrValue := reflect.ValueOf( qptr )
	if qptrValue.Kind() != reflect.Ptr {
		return errors.New( "The argument qptr must be a pointer" )
	}
	qValue := qptrValue.Elem()
	qType := qValue.Type()
	if qType.Kind() != reflect.Struct {
		return errors.New( "The argument qptr must be a pointer to a structure" )
	}
	haveEnqueue := false
	haveDequeue := false
	var elementType reflect.Type = nil
	for i := 0; i != qType.NumField(); i++ {
		field := qType.Field( i )
		tagstring := field.Tag.Get( queue )
		switch tagstring {
		case enqueue:
			if field.Type.Kind() != reflect.Func {
				return fmt.Errorf( "Field '%s' must be a function", field.Name )
			}
			if field.Type.NumIn() != 1 {
				return fmt.Errorf( "Function '%s' must take exactly one argument", field.Name )
			}
			if field.Type.NumOut() != 0 {
				return fmt.Errorf( "Function '%s' must not return anything", field.Name )
			}
			if elementType == nil {
				elementType = field.Type.In( 0 )
			} else {
				if elementType != field.Type.In( 0 ) {
					return fmt.Errorf( "Argument to function '%s' has wrong type '%s', expected '%s'", field.Name, field.Type.In( 0 ).Name(), elementType.Name() )
				}
			}
			qValue.Field( i ).Set( factory.makeEnqueue( field.Type ) )
			haveEnqueue = true
		case dequeue:
			if field.Type.Kind() != reflect.Func {
				return fmt.Errorf( "Field '%s' must be a function", field.Name )
			}
			if field.Type.NumIn() != 0 {
				return fmt.Errorf( "Function '%s' must not take any arguments", field.Name )
			}
			if field.Type.NumOut() != 2 {
				return fmt.Errorf( "Function '%s' must return exactly two values", field.Name )
			}
			if elementType == nil {
				elementType = field.Type.Out( 0 )
			} else {
				if elementType != field.Type.Out( 0 ) {
					return fmt.Errorf( "First return value of function '%s' has wrong type '%s', expected '%s'", field.Name, field.Type.Out( 0 ).Name(), elementType.Name() )
				}
			}
			if field.Type.Out( 1 ).Kind() != reflect.Bool {
				return fmt.Errorf( "Second return value of function '%s' must have type bool", field.Name )
			}
			qValue.Field( i ).Set( factory.makeDequeue( field.Type ) )
			haveDequeue = true
		default:
			continue
		}
	}
	if !haveEnqueue || !haveDequeue {
		return errors.New( "Passed structure must have enqueue and dequeue tags" )
	}
	factory.commit()

	return nil
}
