package fiberaction

import "html/template"

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
	"errors"
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
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
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

type MeetTemp struct {
	ModuleName string          `json:"module_name"`
	Messages   []template.HTML `json:"messages"`
}

var meetTemp = `// DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package {{ .ModuleName }}
{{ range $value := .Messages }}
{{ $value }}
{{ end }}`

type ServiceGroupTemp struct {
	ModName     string
	ModuleName  string
	ServicePath string
	ServicePkg  string
	Name        string
	ParamsPath  string
}

var serviceGroupTemp = `package {{ .ModuleName }}

import (
	"{{ .ModName }}/{{ .ServicePath }}"
	meet "{{ .ModName }}/{{ .ParamsPath }}"
)

type {{ .Name }} struct {
	{{ .ServicePkg }}.Service
}
`

type ServiceFuncTemp struct {
	GroupName string
	Name      string
	Desc      string
}

var serviceFuncTemp = `
// {{ .Name }} {{ .Desc }}
func (s {{ .GroupName }}) {{ .Name }}(req *meet.{{ .GroupName }}{{ .Name }}Req) (*meet.{{ .GroupName }}{{ .Name }}Res, error) {
	//TODO 实现业务
	return &meet.{{ .GroupName }}{{ .Name }}Res{}, nil
}
`

type HandlerControllerGenTemp struct {
	ModuleName    string
	ModuleRoute   string
	HasController bool
	Controllers   []string
	Messages      []template.HTML
	ModName       string
	ServicePath   string
}

var handlerControllerGenTemp = `//DO NOT EDIT.

package {{ .ModuleName }}

import (
	"{{ .ModName }}/internal/fiber/server"
	action "github.com/tdeken/fiberaction"
	"{{ .ModName }}/{{ .ServicePath }}/{{ .ModuleName }}"
	"github.com/gofiber/fiber/v2"
)

// Route 模块路由
func Route() {
	r := server.Web.Server.Group("{{ .ModuleRoute }}")
	{{ if .HasController }}
	action.AutoRegister(r{{ range $value := .Controllers }},
	    {{ $value }}{}{{end}},
    )
	{{ else }}
	action.AutoRegister(r)
	{{ end }}
}
{{ range $value := .Messages }}
{{ $value }}
{{ end }}`

type HandlerGroupTemp struct {
	Name       string
	Desc       string
	Route      string
	Actions    []template.HTML
	ModuleName string
}

var handlerGroupTemp = `// {{ .Name }} {{ .Desc }}
type {{ .Name }} struct {
	Controller
}

// Group 基础请求组
func (s {{ .Name }}) Group() string {
	return "{{ .Route }}"
}

// Register 注册路由
func (s {{ .Name }}) Register() []action.Action {
	return []action.Action{ {{ range $value := .Actions }}
{{ $value }}{{ end }}
	}
}

// 获取依赖服务
func (s {{ .Name }}) getDep(ctx *fiber.Ctx) {{ .ModuleName }}.{{ .Name }} {
	dep := {{ .ModuleName }}.{{ .Name }}{}
	dep.Init(ctx)
	return dep
}`

type HandlerControllerTemp struct {
	ModName     string
	ModuleName  string
	HandlerPath string
	HandlerPkg  string
}

var handlerControllerTemp = `package {{ .ModuleName }}

import (
	"{{ .ModName }}/{{ .HandlerPath }}"
	"github.com/gofiber/fiber/v2"
	action "github.com/tdeken/fiberaction"
)

type Controller struct {
	{{ .HandlerPkg }}.Handler
}

// ChooseMid 可以选择的服务中间件
func (c Controller) ChooseMid(t action.MidType) []fiber.Handler {
	switch t {
	default:
		return nil
	}
}

`

type HandlerFileTemp struct {
	ModuleName string
	ModName    string
	ParamsPath string
}

var handlerFileTemp = `package {{ .ModuleName }}

import (
	"github.com/gofiber/fiber/v2"
	meet "{{ .ModName }}/{{ .ParamsPath }}/{{ .ModuleName }}"
	"{{ .ModName }}/internal/fiber/result"
)
`

type HandlerFuncTemp struct {
	Name   string
	Desc   string
	Group  string
	Req    string
	Route  string
	Method string
}

var handlerFuncTemp = `
// {{ .Name }} {{ .Desc }}
// @Router {{ .Route }} [{{ .Method }}]
func (s {{ .Group }}) {{ .Name }}(ctx *fiber.Ctx) (e error) {
	var form = &meet.{{ .Req }}{}
	if err := s.ValidateRequest(ctx, form); err != nil {
		return result.Json(ctx, nil, err)
	}
	
	res, err := s.getDep(ctx).{{ .Name }}(form)
	return result.Json(ctx, res, err)
}
`

type RouteTemp struct {
	Pkgs    []string
	Modules []string
}

var routeTemp = `//DO NOT EDIT. DO NOT EDIT. DO NOT EDIT.

package route

import ({{ range $value := .Pkgs }}
	"{{ $value }}"{{ end }}
)

func Route() { {{ range $value := .Modules }}
	{{ $value }}.Route(){{ end }}
}
`
