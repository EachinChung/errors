package errors_test

import (
	"fmt"

	"github.com/eachinchung/errors"
)

func ExampleNew() {
	err := errors.New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func ExampleNew_printf() {
	err := errors.New("whoops")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops
	// github.com/eachinchung/errors_test.ExampleNew_printf
	//          /Users/eachin/GolandProjects/errors/example_test.go:17
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWithMessage() {
	cause := errors.New("whoops")
	err := errors.WithMessage(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func ExampleWithStack() {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Println(err)

	// Output: whoops
}

func ExampleWithStack_printf() {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Printf("%+v", err)

	// Example Output:
	// whoops
	// github.com/eachinchung/errors_test.ExampleWithStack_printf
	//          /Users/eachin/GolandProjects/errors/example_test.go:55
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
	// github.com/eachinchung/errors_test.ExampleWithStack_printf
	//          /Users/eachin/GolandProjects/errors/example_test.go:56
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrap() {
	cause := errors.New("whoops")
	err := errors.Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func ExampleCodef() {
	var err error

	err = errors.Codef(1, "this is an error message")
	fmt.Println(err)

	err = errors.Wrap(err, "this is a wrap error message with error code not change")
	fmt.Println(err)

	err = errors.WithCodef(err, 3, "this is a wrap error message with new error code")
	fmt.Println(err)

	// Output:
	// this is an error message
	// this is a wrap error message with error code not change
	// this is a wrap error message with new error code
}

func ExampleCodef_printf() {
	var err error

	err = errors.Codef(1, "this is an error message")
	err = errors.Wrap(err, "this is a wrap error message with error code not change")
	err = errors.WithCodef(err, 3, "this is a wrap error message with new error code")

	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%-v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#-v\n", err)
	fmt.Printf("%#+v\n", err)

	// Example Output:
	// this is a wrap error message with new error code
	// this is a wrap error message with new error code - #2 [/Users/eachin/GolandProjects/errors/example_test.go:128 (github.com/eachinchung/errors_test.ExampleWithCode_printf)] (3) Encoding failed due to an error with the data
	// this is a wrap error message with new error code - #2 [/Users/eachin/GolandProjects/errors/example_test.go:128 (github.com/eachinchung/errors_test.ExampleWithCode_printf)] (3) Encoding failed due to an error with the data; this is a wrap error message with error code not change - #1 [/Users/eachin/GolandProjects/errors/example_test.go:127 (github.com/eachinchung/errors_test.ExampleWithCode_printf)] (1) 内部服务器错误; this is an error message - #0 [/Users/eachin/GolandProjects/errors/example_test.go:126 (github.com/eachinchung/errors_test.ExampleWithCode_printf)] (1) 内部服务器错误
	// [{"error":"Encoding failed due to an error with the data"}]
	// [{"caller":"#2 /Users/eachin/GolandProjects/errors/example_test.go:128 (github.com/eachinchung/errors_test.ExampleWithCode_printf)","code":3,"error":"this is a wrap error message with new error code","message":"Encoding failed due to an error with the data"}]
	// [{"caller":"#2 /Users/eachin/GolandProjects/errors/example_test.go:128 (github.com/eachinchung/errors_test.ExampleWithCode_printf)","code":3,"error":"this is a wrap error message with new error code","message":"Encoding failed due to an error with the data"},{"caller":"#1 /Users/eachin/GolandProjects/errors/example_test.go:127 (github.com/eachinchung/errors_test.ExampleWithCode_printf)","code":1,"error":"this is a wrap error message with error code not change","message":"内部服务器错误"},{"caller":"#0 /Users/eachin/GolandProjects/errors/example_test.go:126 (github.com/eachinchung/errors_test.ExampleWithCode_printf)","code":1,"error":"this is an error message","message":"内部服务器错误"}]
}
func ExampleParseCoder() {
	err := errors.Codef(3, "errors.ParseCoder")
	code := errors.ParseCoder(err)

	fmt.Println(err.Error())
	fmt.Println(code.Code())
	fmt.Println(code.String())
	fmt.Println(code.HTTPStatus())

	// Output:
	// errors.ParseCoder
	// 3
	// encoding failed due to an error with the data
	// 500
}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	return errors.Wrap(e3, "outer")
}

func ExampleCause() {
	err := fn()
	fmt.Println(err)
	fmt.Println(errors.Cause(err))

	// Output: outer: middle: inner: error
	// error
}

func ExampleWrap_extended() {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Example output:
	// error
	// github.com/eachinchung/errors_test.fn
	//          /Users/eachin/GolandProjects/errors/example_test.go:100
	// github.com/eachinchung/errors_test.ExampleWrap_extended
	//          /Users/eachin/GolandProjects/errors/example_test.go:116
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
	// inner
	// github.com/eachinchung/errors_test.fn
	//          /Users/eachin/GolandProjects/errors/example_test.go:101
	// github.com/eachinchung/errors_test.ExampleWrap_extended
	//          /Users/eachin/GolandProjects/errors/example_test.go:116
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
	// middle
	// github.com/eachinchung/errors_test.fn
	//          /Users/eachin/GolandProjects/errors/example_test.go:102
	// github.com/eachinchung/errors_test.ExampleWrap_extended
	//          /Users/eachin/GolandProjects/errors/example_test.go:116
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
	// outer
	// github.com/eachinchung/errors_test.fn
	//          /Users/eachin/GolandProjects/errors/example_test.go:103
	// github.com/eachinchung/errors_test.ExampleWrap_extended
	//          /Users/eachin/GolandProjects/errors/example_test.go:116
	// testing.runExample
	//          /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//          /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//          /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//          /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//          /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
}

func ExampleWrapf() {
	cause := errors.New("whoops")
	err := errors.Wrapf(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Output: oh noes #2: whoops
}

func ExampleErrorf_extended() {
	err := errors.Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Example output:
	// whoops: foo
	// github.com/eachinchung/errors_test.ExampleErrorf_extended
	//         /Users/eachin/GolandProjects/errors/example_test.go:199
	// testing.runExample
	//         /usr/local/opt/go/libexec/src/testing/run_example.go:64
	// testing.runExamples
	//         /usr/local/opt/go/libexec/src/testing/example.go:44
	// testing.(*M).Run
	//         /usr/local/opt/go/libexec/src/testing/testing.go:1505
	// main.main
	//         _testmain.go:91
	// runtime.main
	//         /usr/local/opt/go/libexec/src/runtime/proc.go:255
	// runtime.goexit
	//         /usr/local/opt/go/libexec/src/runtime/asm_amd64.s:1581
}

func Example_stackTrace() {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	err, ok := errors.Cause(fn()).(stackTracer)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	st := err.StackTrace()
	fmt.Printf("%+v", st[0:2]) // top two frames

	// Example output:
	// github.com/eachinchung/errors_test.fn
	//         /Users/eachin/GolandProjects/errors/example_test.go:100
	// github.com/eachinchung/errors_test.ExampleStackTrace
	//         /Users/eachin/GolandProjects/errors/example_test.go:225
}

func ExampleCause_printf() {
	err := errors.Wrap(func() error {
		return func() error {
			return errors.New("hello world")
		}()
	}(), "failed")

	fmt.Printf("%v", err)

	// Output: failed: hello world
}
