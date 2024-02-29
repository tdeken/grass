package fiberweb

import "html/template"

const fiberServerTemp = `package server

import (
	"github.com/gofiber/fiber/v2"
)

type WebServer struct {
	Server *fiber.App
}

var Web = server()

func server() *WebServer {
	ser := fiber.New()

	return &WebServer{
		Server: ser,
	}
}
`

const fiberValidateTemp = `package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

var validate = newValidate()

type RequestInterface any

func CheckParams(ctx *fiber.Ctx, rt RequestInterface) (err error) {
	switch strings.ToUpper(string(ctx.Request().Header.Method())) {
	case "GET":
		err = ctx.QueryParser(rt)
		if err != nil {
			return err
		}
	case "POST", "PUT":
		err = ctx.BodyParser(rt)
		if err != nil {
			return err
		}
	}
	err = validate.Struct(rt)
	return
}

func newValidate() *validator.Validate {
	v := validator.New()

	//load custom validate

	return v
}
`

const fiberRouteTemp = `//DO NOT EDIT.

package route

func Route() { 
	// load handler controller
}
`

type FiberResultTemp struct {
	ModName string
}

const fiberResultTemp = `package result

import (
	"{{ .ModName }}/internal/code"
	"github.com/gofiber/fiber/v2"
)

type WebResult struct {
	Code int32  ` + "`json:" + `"code"` + "`" + ` //错误码
	Msg  string ` + "`json:" + `"msg"` + "`" + `  //返回的消息
	Data any    ` + "`json:" + `"data"` + "`" + ` //返回的数据结果
}

func Json(ctx *fiber.Ctx, res any, err error) error {
	var crr *code.Error
	if err != nil {
		var ok bool
		crr, ok = err.(*code.Error)
		if !ok {
			crr = code.NewError(code.UndefinedErrorCode, err.Error())
		}
	}

	var data = res
	if data == nil {
		data = make(map[string]interface{})
	}

	return ctx.JSON(WebResult{
		Code: crr.GetCode(),
		Msg:  crr.GetDetail(),
		Data: data,
	})
}
`

type FiberTemp struct {
	ModName string
}

const fiberTemp = `package fiber

import (
	"context"
	"{{ .ModName }}/internal/config"
	"{{ .ModName }}/internal/fiber/route"
	"{{ .ModName }}/internal/fiber/server"
	"fmt"
	"time"
)

func Run() (err error) {
	//注册路由
	route.Route()
	
	// do custom route
	
	//启动项目
	if err = server.Web.Server.Listen(fmt.Sprintf(":%d", config.Conf.Server.Port)); err != nil {
		return
	}
	return
}

func Shutdown() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Web.Server.ShutdownWithContext(ctx)
}
`

const codeTemp = `package code

const OK = 200
const OKMsg = "ok"


const (
	SystemErrorCode    = -1 //系统级别错误
	UndefinedErrorCode = 0  //没有定义错误类型的都用这个
)
`

const codeErrorTemp = `package code

import (
	"encoding/json"
)

type Error struct {
	Form   string ` + "`json:" + `"form"` + "`" + `  //来源
	Code   int32  ` + "`json:" + `"code"` + "`" + `   //错误码
	Detail string ` + "`json:" + `"detail"` + "`" + ` //错误信息
}

// GetCode 状态码
func (e *Error) GetCode() int32 {
	if e == nil {
		return OK
	}
	return e.Code
}

// GetDetail 状态码说明
func (e *Error) GetDetail() string {
	if e == nil {
		return OKMsg
	}
	return e.Detail
}

func (e *Error) Error() string {
	if e == nil {
		return "nil"
	}
	b, _ := json.Marshal(e)
	return string(b)
}

// NewError 实例化一个错误
func NewError(code int32, detail string) *Error {
	return &Error{Code: code, Detail: detail}
}

// NewFormError 实例化一个有来源的错误
func NewFormError(code int32, detail, form string) *Error {
	return &Error{Form: form, Code: code, Detail: detail}
}
`

const confingTemp = `package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf = &Config{}

// FilePath 配置文件路径信息
type FilePath struct {
	ConfigName string // 设置配置文件名 (不需要后缀)
	ConfigType string // 设置配置文件类型
	ConfigPath string // 设置配置文件路径
}

// LoadConfig 从文件中加载配置
func LoadConfig(filepath FilePath) error {
	viper.SetConfigType(filepath.ConfigType) // 设置配置文件类型
	viper.SetConfigName(filepath.ConfigName) // 设置配置文件名 (不需要后缀)
	viper.AddConfigPath(filepath.ConfigPath) // 设置配置文件路径

	err := viper.ReadInConfig() // 读取配置文件
	if err != nil {
		return fmt.Errorf("failed to read etc file: %s", err)
	}

	err = viper.Unmarshal(Conf) // 将读取的配置映射到 Config 结构体中
	if err != nil {
		return fmt.Errorf("failed to unmarshal etc file: %s", err)
	}

	return nil
}

// Config 本地配置
type Config struct {
	Server Server      ` + "`mapstructure:" + `"server"` + "`" + ` //系统配置
}

type Server struct {
	Env         string ` + "`mapstructure:" + `"env"` + "`" + `          // 动环境，取值为 "dev"、"test" 或 "prod"，默认为 "dev"
	Port        int    ` + "`mapstructure:" + `"port"` + "`" + `         // 服务端口
}
`

type BootTemp struct {
	ModName string
}

const bootTemp = `package boot

import (
	"{{ .ModName }}/internal/app"
	"{{ .ModName }}/internal/config"
	"{{ .ModName }}/internal/fiber"
	"flag"
	"fmt"
	"log"
)

// Init 初始化项目配置
func Init() {
	var env string
	flag.StringVar(&env, "env", "local", "启动环境")
	flag.Parse()

	err := config.LoadConfig(config.FilePath{
		ConfigName: fmt.Sprintf("config-%s", env),
		ConfigType: "yaml",
		ConfigPath: "etc",
	})
	if err != nil {
		log.Fatalf("配置文件加载失败:%v", err)
	}

}

func Run() {
	go fiber.Run()

	return
}

// Shutdown 关闭运行程序
func Shutdown() {
	if err := fiber.Shutdown(); err != nil {
		panic(err)
	}
}
`

type MainTemp struct {
	ModName string
	Spc     template.HTML
}

const mainTemp = `package main

import (
	"{{ .ModName }}/internal/boot"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//加载配置，系统全局变量
	boot.Init()

	//启动服务
	boot.Run()

	//监听程序退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	select {
	case {{ .Spc }}ch:
		defer close(ch)
	}

	//关闭程序
	boot.Shutdown()
}
`
