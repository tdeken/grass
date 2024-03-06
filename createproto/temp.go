package createproto

type ProtoTemp struct {
	Route string
}

const protoTemp = `#路由组
route: "{{ .Route }}"

### swagger doc support tip

#标题
title: "{{ .Route }}接口"
#备注说明
desc: "文档 固定返回格式 {\"code\": 0, \"msg\": \"ok\", data: null } code: 错误码 0为成功, 其他都属于错误类型; msg: 错误信息, 错误原因; data: 返回的数据格式，具体看接口返回"
#支持协议
schemes:
  - "http"
  - "https"
#请求host
host: "127.0.0.1:8080"
#接口版本
ver: "1.0"
#身份校验
auth:
  security: "apiKey"
  title: "BearerToken"
  in: "header"
  name: "Authorization"
  token: ""
#联系方式
contact:
  name: ""
  url: ""
  email: ""`
