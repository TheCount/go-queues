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

// Package queue provides implementations of the queue abstract data type
// based on a simple configuration scheme in a type-safe manner.
// For created queues, type safety is ensured at compile time.
//
// Clients provide their own queue data structure
// modelled after the GenericQueue type.
// The dummy type T in GenericQueue
// should be replaced by the appropriate type.
// A pointer to an instance can then be passed to the Make function,
// which will fill in the queue methods.
// Queue methods are identified by structure tags (see GenericQueue).
// Other structure fields are ignored by the Make function.
//
// A configuration can be passed to the Make function,
// choosing an implementation with specific characteristics for the queue.
// For example, configurable parameters are initial queue buffer size or
// whether the queue should be safe to access concurrently.
// The default configuration yields a queue suitable for most uses.
package queue

// Version information for the queue package.
const(
	VersionMajor = 0 // remains at 0 for now
	VersionMinor = 0 // remains at 0 while in alpha
	VersionMicro = 0 // increased with every release
)
