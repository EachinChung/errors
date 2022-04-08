package errors

//goland:noinspection SpellCheckingInspection
import (
	pkgerrors "github.com/pkg/errors"
)

// New 返回带有消息的错误
// New 在它被调用的地方，记录堆栈跟踪
func New(message string) error {
	return pkgerrors.New(message)
}

// Errorf 根据格式说明符进行格式化，并将字符串作为 error 的值返回
// Errorf 在它被调用的地方，记录堆栈跟踪
func Errorf(format string, args ...interface{}) error {
	return pkgerrors.Errorf(format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	return pkgerrors.WithStack(err)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return pkgerrors.Wrap(err, message)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return pkgerrors.Wrapf(err, format, args...)
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	return pkgerrors.WithMessage(err, message)
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	return pkgerrors.WithMessagef(err, format, args...)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	return pkgerrors.Cause(err)
}
