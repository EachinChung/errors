package errors

//goland:noinspection SpellCheckingInspection
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// New 返回带有消息的 error
// New 在它被调用的地方，记录堆栈跟踪
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf 根据格式说明符进行格式化，并将字符串作为 error 的值返回
// Errorf 在它被调用的地方，记录堆栈跟踪
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// WithStack 在调用堆栈时，用堆栈跟踪注释 error
// 如果 err 为 nil，WithStack 返回 nil
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			msg:   e.msg,
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	return &withStack{err, callers()}
}

// Wrap 返回 error，该错误用 Wrap 堆栈跟踪注释 err，并返回提供错误信息
// 如果 err 为 nil，则 Wrap 返回 nil
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			msg:   message,
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{cause: err, msg: message}
	return &withStack{err, callers()}
}

// Wrapf 返回 error，该错误用 Wrapf 堆栈跟踪注释 err，并返回格式化错误信息
// 如果 err 为 nil，则 Wrapf 返回 nil
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			msg:   fmt.Sprintf(format, args...),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

// WithMessage 用 message 注释错误
// 如果 err 为 nil，则 WithMessage 返回 nil
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// WithMessagef 用格式化的 message 注释错误
// 如果 err 为 nil，则 WithMessagef 返回 nil
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

// Code 用 message 注释错误，同时用错误码映射错误
func Code(code int, message string) error {
	return &withCode{
		msg:   message,
		code:  code,
		stack: callers(),
	}
}

// Codef 用格式化的 message 注释错误，同时用错误码映射错误
func Codef(code int, format string, args ...interface{}) error {
	return &withCode{
		msg:   fmt.Sprintf(format, args...),
		code:  code,
		stack: callers(),
	}
}

// WithCode 返回 error，该错误用 WithCodef 堆栈跟踪注释 err
// 同时用错误码映射错误，并返回错误信息
// 如果 err 为 nil，则 WithCodef 返回 nil
func WithCode(err error, code int, message string) error {
	if err == nil {
		return nil
	}

	return &withCode{
		msg:   message,
		code:  code,
		cause: err,
		stack: callers(),
	}
}

// WithCodef 返回 error，该错误用 WithCodef 堆栈跟踪注释 err
// 同时用错误码映射错误，并返回格式化错误信息
// 如果 err 为 nil，则 WithCodef 返回 nil
func WithCodef(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		msg:   fmt.Sprintf(format, args...),
		code:  code,
		cause: err,
		stack: callers(),
	}
}

// fundamental 一个错误，它有一个 msg 和 stack，但没有调用者
type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

//nolint:errcheck
//goland:noinspection GoUnhandledErrorResult
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

// Unwrap 提供 Go 1.13 错误链的兼容性
func (w *withStack) Unwrap() error { return w.error }

//nolint:errcheck
//goland:noinspection GoUnhandledErrorResult
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }

// Unwrap 提供 Go 1.13 错误链的兼容性
func (w *withMessage) Unwrap() error { return w.cause }

//nolint:errcheck
//goland:noinspection GoUnhandledErrorResult
func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

type withCode struct {
	msg   string
	code  int
	cause error
	*stack
}

func (w *withCode) Error() string { return w.msg }

// Cause 返回 error 的原因
func (w *withCode) Cause() error { return w.cause }

// Unwrap 提供 Go 1.13 兼容性
func (w *withCode) Unwrap() error { return w.cause }

// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Verbs:
//     %s  -   如果没有指定，则返回映射到错误代码或错误消息的用户安全错误字符串。
//     %v      %s 的别名
//
// Flags:
//      #      JSON 格式的输出，用于日志记录
//      -      输出调用者详细信息，有助于故障排除
//      +      输出完整的错误堆栈详细信息，对调试有用
//goland:noinspection GoUnhandledErrorResult
func (w *withCode) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})
		var jsonData []map[string]interface{}

		var (
			flagDetail bool
			flagTrace  bool
			modeJSON   bool
		)

		if state.Flag('#') {
			modeJSON = true
		}

		if state.Flag('-') {
			flagDetail = true
		}
		if state.Flag('+') {
			flagTrace = true
		}

		sep := ""
		errs := list(w)
		length := len(errs)
		for k, e := range errs {
			info := buildFormatInfo(e)
			jsonData, str = format(length-k-1, jsonData, str, info, sep, flagDetail, flagTrace, modeJSON)
			sep = "; "

			if !flagTrace {
				break
			}
		}
		if modeJSON {
			var b []byte
			b, _ = json.Marshal(jsonData)

			str.Write(b)
		}

		fmt.Fprintf(state, "%s", strings.Trim(str.String(), "\r\n\t"))
	default:
		fmt.Fprintf(state, buildFormatInfo(w).err)
	}
}

// Cause 如果可能的话，返回 error 的根本原因
// 如果实现以下接口，则 error 会返回原因:
//
//     type causer interface {
//            Cause() error
//     }
//
// 如果 error 没有实现 Cause，则返回原始 error。
// 如果 error 为 nil，则将返回 nil，而无需进一步调查。
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}

		if cause.Cause() == nil {
			break
		}

		err = cause.Cause()
	}
	return err
}
