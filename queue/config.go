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

// Flags is a bitflag type to hold information about queue configuration.
type Flags uint64

// These flags determine some aspects of queue configurations,
// such as access patterns.
const(
	// FNonConcurrent indicates that
	// the queue must not be accessed concurrently.
	// This flag is mutually exclusive with both FMultiReader and FMultiWriter
	FNonConcurrent Flags = 1 << iota

	// FMultiReader indicates that multiple goroutines
	// can dequeue safely from the queue concurrently.
	// This flag is mutually exclusive with FNonConcurrent
	FMultiReader

	// FMultiWriter indicates that multiple goroutines
	// can enqueue safely to the queue concurrently.
	// This flag is mutually exclusive with FNonConcurrent
	FMultiWriter

	// FNotImplemented forces the configurator to assume that there is no
	// implementation for the specified configuration.
	// Can be used for testing purposes.
	FNotImplemented Flags = 1 << 63
)

// DefaultInitialCapacity is the initial capacity used by DefaultConfig()
const DefaultInitialCapacity = 4

// Config holds the configuration for a queue
type Config struct {
	// Flags denotes the configuration flags for Make.
	Flags Flags

	// initialCapacity denotes the initial capacity of the queue.
	initialCapacity int
}

// IsValid checks whether the configuration is valid.
func ( c *Config ) IsValid() bool {
	if ( ( c.Flags & FNonConcurrent ) != 0 ) && ( ( c.Flags & ( FMultiReader | FMultiWriter ) ) != 0 ) {
		return false
	} else {
		return true
	}
}

// NonConcurrent selects queue types whose instances are safe to access
// only sequentially.
func ( c *Config ) NonConcurrent() *Config {
	c.Flags &= ^( FMultiReader | FMultiWriter )
	c.Flags |= FNonConcurrent

	return c
}

// SingleReader selects queue types whose Dequeue method may not be called
// concurrently.
func ( c *Config ) SingleReader() *Config {
	c.Flags &= ^FMultiReader

	return c
}

// MultiReader selects queue types whose Dequeue method may be called
// concurrently from multiple goroutines.
func ( c *Config ) MultiReader() *Config {
	c.Flags &= ^FNonConcurrent
	c.Flags |= FMultiReader

	return c
}

// SingleWriter selects queue types whose Enqueue method may be not called
// concurrently.
func ( c *Config ) SingleWriter() *Config {
	c.Flags &= ^FMultiWriter

	return c
}

// MultiWriter selects queue types whose Enqueue method may be called
// concurrently from multiple goroutines.
func ( c *Config ) MultiWriter() *Config {
	c.Flags &= ^FNonConcurrent
	c.Flags |= FMultiWriter

	return c
}

// InitialCapacity sets the initial capacity for the queue.
// A negative value or a very small non-negative value will be increased
// to the minimum capacity for the selected queue automatically.
func ( c *Config ) InitialCapacity( capacity int ) *Config {
	c.initialCapacity = capacity

	return c
}

// DefaultConfig returns a default configuration suitable for most uses.
func DefaultConfig() *Config {
	return &Config{
		Flags: FMultiReader | FMultiWriter,
		initialCapacity: DefaultInitialCapacity,
	}
}

// factory returns a factory for this configuration.
// If no matching implementation exists, nil is returned.
func ( c *Config ) factory() factory {
	if ( c.Flags & FNotImplemented ) != 0 {
		return nil
	}
	if ( c.Flags & FNonConcurrent ) != 0 {
		return newSimpleQueueFactory( c.initialCapacity )
	} else {
		return newLockedQueueFactory( c.initialCapacity )
	}
}
