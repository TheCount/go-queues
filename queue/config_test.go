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
	"testing"
)

func TestIsValid( t *testing.T ) {
	config := DefaultConfig()
	if !config.IsValid() {
		t.Error( "Default configuration is not valid" )
	}
	config.Flags = 0
	if !config.IsValid() {
		t.Error( "Single reader/single writer config is not valid" )
	}
	config.Flags = FNonConcurrent
	if !config.IsValid() {
		t.Error( "Non-concurrent configuration is not valid" )
	}
	config.Flags = FMultiReader
	if !config.IsValid() {
		t.Error( "Multi reader/single writer config is not valid" )
	}
	config.Flags = FNonConcurrent | FMultiReader
	if config.IsValid() {
		t.Error( "Multi reader vs. non-concurrent config is valid" )
	}
	config.Flags = FMultiWriter
	if !config.IsValid() {
		t.Error( "Single reader/multi writer config is not valid" )
	}
	config.Flags = FNonConcurrent | FMultiWriter
	if config.IsValid() {
		t.Error( "Multi writer vs. non-concurrent config is valid" )
	}
	config.Flags = FMultiReader | FMultiWriter
	if !config.IsValid() {
		t.Error( "Multi reader/multi writer config is not valid" )
	}
	config.Flags = FNotImplemented
	if !config.IsValid() {
		t.Error( "Not implemented config is not valid" )
	}
}

func TestDefault( t *testing.T ) {
	config := DefaultConfig()
	if config.initialCapacity != DefaultInitialCapacity {
		t.Errorf( "Default configuration capacity %d != default initial capacity %d", config.initialCapacity, DefaultInitialCapacity )
	}
}

func TestNonConcurrent( t *testing.T ) {
	config := DefaultConfig()
	config.NonConcurrent()
	if ( config.Flags & FNonConcurrent ) == 0 {
		t.Error( "Non-concurrent flag not set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after NonConcurrent()" )
	}
	config.Flags &= ^FNonConcurrent
	config.NonConcurrent()
	if ( config.Flags & FNonConcurrent ) == 0 {
		t.Error( "Non-concurrent flag not set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after NonConcurrent()" )
	}
}

func TestSingleReader( t *testing.T ) {
	config := DefaultConfig()
	config.SingleReader()
	if ( config.Flags & FMultiReader ) != 0 {
		t.Error( "Multi reader flag set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after SingleReader()" )
	}
	config.Flags |= FMultiReader
	config.SingleReader()
	if ( config.Flags & FMultiReader ) != 0 {
		t.Error( "Multi reader flag set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after SingleReader()" )
	}
}

func TestMultiReader( t *testing.T ) {
	config := DefaultConfig()
	config.MultiReader()
	if ( config.Flags & FMultiReader ) == 0 {
		t.Error( "Multi reader flag not set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after MultiReader()" )
	}
	config.Flags &= ^FMultiReader
	config.MultiReader()
	if ( config.Flags & FMultiReader ) == 0 {
		t.Error( "Multi reader flag not set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after MultiReader()" )
	}
}

func TestSingleWriter( t *testing.T ) {
	config := DefaultConfig()
	config.SingleWriter()
	if ( config.Flags & FMultiWriter ) != 0 {
		t.Error( "Multi writer flag set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after SingleWriter()" )
	}
	config.Flags |= FMultiWriter
	config.SingleWriter()
	if ( config.Flags & FMultiWriter ) != 0 {
		t.Error( "Multi writer flag set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after SingleWriter()" )
	}
}

func TestMultiWriter( t *testing.T ) {
	config := DefaultConfig()
	config.MultiWriter()
	if ( config.Flags & FMultiWriter ) == 0 {
		t.Error( "Multi writer flag not set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after MultiWriter()" )
	}
	config.Flags &= ^FMultiWriter
	config.MultiWriter()
	if ( config.Flags & FMultiWriter ) == 0 {
		t.Error( "Multi writer flag not set" )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after MultiWriter()" )
	}
}

func TestConfigCapacity( t *testing.T ) {
	config := DefaultConfig()
	config.InitialCapacity( 42 )
	if config.initialCapacity != 42 {
		t.Errorf( "Bad capacity: %d, expected 42", config.initialCapacity )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after InitialCapacity()" )
	}
	config.InitialCapacity( -42 )
	if config.initialCapacity > DefaultInitialCapacity {
		t.Errorf( "Bad capacity choice after InitialCapacity() with negative argument: %d", config.initialCapacity )
	}
	if !config.IsValid() {
		t.Error( "Configuration not valid after InitialCapacity()" )
	}
}

func TestFactory( t *testing.T ) {
	config := DefaultConfig()
	factory := config.factory()
	if factory == nil {
		t.Error( "Default configuration not implemented" )
	}
	config.Flags |= FNotImplemented
	factory = config.factory()
	if factory != nil {
		t.Error( "Not implemented configuration implemented" )
	}
}
