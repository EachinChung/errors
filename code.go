package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder defaultCoder = defaultCoder{1, http.StatusInternalServerError, "内部服务器错误"}
)

// Coder 定义错误代码详细信息的接口。
type Coder interface {
	// HTTPStatus 错误对应的HTTP状态码
	HTTPStatus() int

	// String 外部 (用户) 面对的错误信息
	String() string

	// Code 返回错误码
	Code() int
}

type defaultCoder struct {
	// C 错误码
	C int

	// HTTP 关联的HTTP状态码。
	HTTP int

	// Ext 外部 (用户) 面对错误文本。
	Ext string
}

// Code 返回错误码
func (coder defaultCoder) Code() int {
	return coder.C
}

// String 外部 (用户) 面对的错误信息
func (coder defaultCoder) String() string {
	return coder.Ext
}

// HTTPStatus 返回关联的HTTP状态代码 (如果有)。否则，返回200。
func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return 500
	}

	return coder.HTTP
}

// codes 包含错误代码到元数据的映射。
var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// Register 注册用户定义错误代码。
func Register(coder Coder) {
	if 0 <= coder.Code() && coder.Code() <= 100 {
		panic("code '0 ~ 100' is the reserved error code of the package `github.com/eachinchung/errors`")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister 注册用户定义错误代码。
// 当已存在相同的 code 时，会发生 panic
func MustRegister(coder Coder) {
	if 0 <= coder.Code() && coder.Code() <= 100 {
		panic("code '0 ~ 100' is the reserved error code of the package `github.com/eachinchung/errors`")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder 将任何错误解析为 *withCode。
// nil 错误将直接返回 nil。
// 没有 withStack 的错误，将被解析为 ErrUnknown.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

// IsCode 报告err链中的任何错误是否包含给定的错误代码。
func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}
