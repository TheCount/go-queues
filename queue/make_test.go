package queue

import(
	"testing"
)

type structEmpty struct {
}

type structGarbage struct {
	foo int
	bar float64
}

type structMissingTags struct {
	Enqueue func( int )
	Dequeue func() ( int, bool )
}

type structBadTags struct {
	Enqueue func( int ) `enqueue`
	Dequeue func() ( int, bool ) `dequeue`
}

type structUnhandledTags struct {
	Enqueue func( int ) `queue:"foo"`
	Dequeue func() ( int, bool ) `queue:"bar"`
}

type structEnqueueNotFunc struct {
	Enqueue *func( int ) `queue:"enqueue"`
	Dequeue func() ( int, bool ) `queue:"dequeue"`
}

type structDequeueNotFunc struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue *func() ( int, bool ) `queue:"dequeue"`
}

type structMissingEnqueue struct {
	Dequeue func() ( int, bool ) `queue:"dequeue"`
}

type structMissingDequeue struct {
	Enqueue func( int ) `queue:"enqueue"`
}

type structBadEnqueue1 struct {
	Enqueue func() `queue:"enqueue"`
	Dequeue func() ( int, bool ) `queue:"dequeue"`
}

type structBadEnqueue2 struct {
	Enqueue func( int, bool ) `queue:"enqueue"`
	Dequeue func() ( int, bool ) `queue:"dequeue"`
}

type structBadEnqueue3 struct {
	Enqueue func( int ) int `queue:"enqueue"`
	Dequeue func() ( int, bool ) `queue:"dequeue"`
}

type structBadDequeue1 struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func() `queue:"dequeue"`
}

type structBadDequeue2 struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func() int `queue:"dequeue"`
}

type structBadDequeue3 struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func() ( int, int ) `queue:"dequeue"`
}

type structBadDequeue4 struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func() ( int, bool, float64 ) `queue:"dequeue"`
}

type structBadDequeue5 struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func( int ) ( int, bool ) `queue:"dequeue"`
}

type structTypeMismatch1 struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func()( float64, bool ) `queue:"dequeue"`
}

type structTypeMismatch2 struct {
	Dequeue func()( float64, bool ) `queue:"dequeue"`
	Enqueue func( int ) `queue:"enqueue"`
}

type structOK struct {
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func() ( int, bool ) `queue:"dequeue"`
}

type structOKExtraFields struct {
	foo int
	Enqueue func( int ) `queue:"enqueue"`
	Dequeue func() ( int, bool ) `queue:"dequeue"`
	bar float64
}

type structOKUnusualNames struct {
	Foo func( int ) `queue:"enqueue"`
	Bar func() ( int, bool ) `queue:"dequeue"`
}

func TestMake( t *testing.T ) {
	// Create configurations
	config := DefaultConfig().NonConcurrent()
	configNotImplemented := DefaultConfig()
	configNotImplemented.Flags |= FNotImplemented
	configInvalid := DefaultConfig()
	configInvalid.Flags |= FNonConcurrent | FMultiReader | FMultiWriter
	// nil interface
	err := Make( nil, nil )
	if err == nil {
		t.Error( "Make succeeded with nil interface" )
	}
	err = Make( nil, config )
	if err == nil {
		t.Error( "Make succeeded with nil interface and non-nil config" )
	}
	// non-pointer type
	err = Make( 42, config )
	if err == nil {
		t.Error( "Make succeeded with argument 42" )
	}
	// pointer to non-struct
	v := 42
	err = Make( &v, config )
	if err == nil {
		t.Error( "Make succeeded with pointer-to-int argument" )
	}
	// Various ill-formed structs
	var se structEmpty
	err = Make( &se, config )
	if err == nil {
		t.Error( "Make succeeded despite empty struct" )
	}
	var sg structGarbage
	err = Make( &sg, config )
	if err == nil {
		t.Error( "Make succeeded despite garbage struct" )
	}
	var smt structMissingTags
	err = Make( &smt, config )
	if err == nil {
		t.Error( "Make succeeded despite missing tags" )
	}
	var sbt structBadTags
	err = Make( &sbt, config )
	if err == nil {
		t.Error( "Make succeeded despite struct with bad tags" )
	}
	var sut structUnhandledTags
	err = Make( &sut, config )
	if err == nil {
		t.Error( "Make succeeded despite struct with unhandled tags" )
	}
	var senf structEnqueueNotFunc
	err = Make( &senf, config )
	if err == nil {
		t.Error( "Make succeeded despite enqueue not being a function" )
	}
	var sdnf structDequeueNotFunc
	err = Make( &sdnf, config )
	if err == nil {
		t.Error( "Make succeeded despite dequeue not being a function" )
	}
	var sme structMissingEnqueue
	err = Make( &sme, config )
	if err == nil {
		t.Error( "Make succeeded despite missing enqueue" )
	}
	var smd structMissingDequeue
	err = Make( &smd, config )
	if err == nil {
		t.Error( "Make succeeded despite missing dequeue" )
	}
	var sbe1 structBadEnqueue1
	err = Make( &sbe1, config )
	if err == nil {
		t.Error( "Make succeeded despite enqueue taking no arguments" )
	}
	var sbe2 structBadEnqueue2
	err = Make( &sbe2, config )
	if err == nil {
		t.Error( "Make succeeded despite enqueue taking too many arguments" )
	}
	var sbe3 structBadEnqueue3
	err = Make( &sbe3, config )
	if err == nil {
		t.Error( "Make succeeded despite enqueue returning a value" )
	}
	var sbd1 structBadDequeue1
	err = Make( &sbd1, config )
	if err == nil {
		t.Error( "Make succeeded despite dequeue returning no value" )
	}
	var sbd2 structBadDequeue2
	err = Make( &sbd2, config )
	if err == nil {
		t.Error( "Make succeeded despite dequeue returning only one value" )
	}
	var sbd3 structBadDequeue3
	err = Make( &sbd3, config )
	if err == nil {
		t.Error( "Make succeeded despite dequeue returning second value not of type bool" )
	}
	var sbd4 structBadDequeue4
	err = Make( &sbd4, config )
	if err == nil {
		t.Error( "Make succeeded despite dequeue returning too many values" )
	}
	var sbd5 structBadDequeue5
	err = Make( &sbd5, config )
	if err == nil {
		t.Error( "Make succeeded despite dequeue taking an argument" )
	}
	var stm1 structTypeMismatch1
	err = Make( &stm1, config )
	if err == nil {
		t.Error( "Make succeeded despite element type mismatch between enqueue and dequeue" )
	}
	var stm2 structTypeMismatch2
	err = Make( &stm2, config )
	if err == nil {
		t.Error( "Make succeeded despite element type mismatch between dequeue and enqueue" )
	}
	// OK structs
	var sok structOK
	err = Make( &sok, config )
	if err != nil {
		t.Error( "Creation of non-concurrent queue failed" )
	}
	err = Make( &sok, configNotImplemented )
	if err == nil {
		t.Error( "Make succeeded despite not implemented configuration" )
	}
	err = Make( &sok, configInvalid )
	if err == nil {
		t.Error( "Make succeeded with invalid configuration" )
	}
	// FIXME: also test with nil config
	var extra structOKExtraFields
	err = Make( &extra, config )
	if err != nil {
		t.Error( "Creation of non-concurrent queue with extra fields failed" )
	}
	var unusual structOKUnusualNames
	err = Make( &unusual, config )
	if err != nil {
		t.Error( "Creation of non-concurrent queue with unusually-named methods failed" )
	}
	var generic GenericQueue
	err = Make( &generic, config )
	if err != nil {
		t.Error( "Creation of generic non-concurrent queue failed" )
	}
}
