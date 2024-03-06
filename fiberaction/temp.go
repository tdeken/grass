package fiberaction

type ServiceTemp struct {
	Pkg     string
	Service string
}

const serviceTemp = `package {{ .Pkg }}

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type {{ .Service }} struct {
	Ctx *fiber.Ctx
}

func (s *{{ .Service }}) Init(ctx *fiber.Ctx) {
	s.Ctx = ctx
}

func (s *{{ .Service }}) Context() *fasthttp.RequestCtx {
	return s.Ctx.Context()
}`

type HandlerTemp struct {
	ModName string
	Pkg     string
	Handler string
}

const handlerTemp = `package {{ .Pkg }}

import (
	"{{ .ModName }}/internal/code"
	"{{ .ModName }}/internal/fiber/validate"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	action "github.com/tdeken/fiberaction"
)

type {{ .Handler }} struct {
}

// ValidateRequest 统一校验请求数据
func (s {{ .Handler }}) ValidateRequest(ctx *fiber.Ctx, rt validate.RequestInterface) *code.Error {
	err := validate.CheckParams(ctx, rt)
	if err != nil {
		var errMsg string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, er := range err.(validator.ValidationErrors) {
				errMsg = fmt.Sprintf("错误字段:%v,验证类型:%v:%v,参数值:%v", er.Field(), er.Tag(), er.Param(), er.Value())
				break
			}
		} else {
			errMsg = err.Error()
		}
		return code.NewError(code.VerifyErrorCode, errMsg)
	}

	return nil
}

// ChooseMid 可以选择的服务中间件
func (s {{ .Handler }}) ChooseMid(t action.MidType) (ms []fiber.Handler) {
	if t == nil {
		return
	}

	return
}`
