//go:build go1.13
// +build go1.13

package errors

//goland:noinspection SpellCheckingInspection
import (
	stderrors "errors"
)

// Is 报告错误链中的任何错误是否与目标匹配
//
// 该链由 error 组成, 然后是通过反复调用 Unwrap 获得的 error 序列。
//
// 如果一个 error 等于 target, 或者如果它实现了一个方法 is (error) bool 使得 is (target) 返回 true,
// 则认为该 error 与该 target 匹配。
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As 找到与目标匹配的错误链中的第一个 error, 如果是, 则将 target 设置为该 error 并返回 true。
//
// 该链由 error 组成, 然后是通过反复调用 Unwrap 获得的 error 序列。
//
// 如果 error 的具体值可分配给 target 指向的值, 或者如果 error 具有方法 As (interface{}) bool,
// 则 error 与 target 匹配, As(target) 返回 true。在后一种情况下, As 方法负责设定 target。
//
// 如果 target 不是指向实现 error 的类型或任何接口类型的非nil指针,
// As 将会 panic。如果错误为 nil, As 返回 false。
//goland:noinspection GoErrorsAs
func As(err error, target interface{}) bool { return stderrors.As(err, target) }

// Unwrap 如果错误的类型包含一个 Unwrap 方法返回错误, 则返回该错误上的 Unwrap 方法的结果。
// 否则, Unwrap 返回 nil
func Unwrap(err error) error { return stderrors.Unwrap(err) }
