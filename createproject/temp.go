package createproject

type ProtoYamlFile struct {
	ModName string `json:"mod_name"`
}

const protoYamlFile = `## NOT EDIT; NOT EDIT; NOT EDIT

mod_name: {{ .ModName }}

source: proto

analyze:
  sources: api/http
  handler: internal/handler
  service: internal/service

swagger:
  path: api/swagger
  code: code
  msg: msg
  data: data`
