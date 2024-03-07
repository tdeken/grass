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

const exampleTemp = `{
  "group": {
    "name": "Example",
    "desc": "示例接口组"
  },
  "interfaces": [
    {
      "name": "Example",
      "desc": "示例接口",
      "method": "POST",
      "mid_type": "mid_key",
      "msgs": [
        {
          "name": "MsgsExample",
          "fields": [
            {"name": "example_1", "class": "string", "desc": "oss accessKeyId"},
            {"name": "example_2", "class": "int64", "desc": "oss accessKeyId"},
            {"name": "example_3", "class": "[]string", "desc": "oss accessKeyId"}
          ]
        },
        {
          "name": "MsgsExample2",
          "fields": [
            {"name": "example_1", "class": "string", "desc": "oss accessKeyId"},
            {"name": "example_2", "class": "[]string", "desc": "oss accessKeyId"},
            {"name": "example_3", "class": "[]~MsgsExample", "desc": "oss accessKeyId"},
            {"name": "example_4", "class": "~MsgsExample", "desc": "oss accessKeyId"}
          ]
        }
      ],
      "req": {
        "fields": [
          {"name": "id", "class": "string", "desc": "id", "validate": "required"}
        ]
      },
      "res": {
        "fields": [
          {"name": "msgs_example", "class": "[]~MsgsExample2", "desc": "msgs_example"},
          {"name": "name", "class": "string", "desc": "name"},
          {"name": "age", "class": "float64", "desc": "age"},
          {"name": "games", "class": "[]?Game", "desc": "games"},
          {"name": "next_play", "class": "?Game", "desc": "games"},
          {"name": "play", "class": "&Game", "desc": "games"}
        ],
        "messages": [
          {
            "name": "Game",
            "fields": [
              {"name": "name", "class": "string", "desc": "name"},
              {"name": "time", "class": "int32", "desc": "time"}
            ]
          }
        ]
      }
    }
  ]
}`

const exampleTomlTemp = `[group]
name = "Example"
desc = "示例接口组"

[[interfaces]]
name = "Example"
desc = "示例接口"
method = "POST"
mid_type = "mid_key"
[[interfaces.msgs]]
name = "MsgsExample"
fields = [
    "n=example_1;c=string;d=example_1",
    "n=example_2;c=string;d=example_2",
    "n=example_3;c=string;d=example_3",
]
[[interfaces.msgs]]
name = "MsgsExample2"
fields = [
    "n=example_1;c=string;d=example_1",
    "n=example_2;c=string;d=example_2",
    "n=example_3;c=string;d=example_3",
]

[interfaces.req]
fileds = [
    "n=id;c=string;d=id;v=required",
]

[interfaces.res]
fileds = [
    "n=example_1;c=[]~MsgsExample2;d=msgs_example",
    "n=name;c=string;d=name",
    "n=age;c=float64;d=age",
    "n=games;c=[]?Game;d=games",
    "n=next_play;c=?Game;d=next_play",
    "n=play;c=&Game;d=play",
]
[[interfaces.res.messages]]
name="Game"
fields = [
    "n=name;c=string;d=name",
    "n=time;c=int32;d=time",
]
`
