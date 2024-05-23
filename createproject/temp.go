package createproject

type ProtoYamlFile struct {
	ModName string `json:"mod_name"`
}

const protoYamlFile = `mod_name: {{ .ModName }}

proto: 
  path: proto
  file_type: json # json toml

analyze:
  sources: internal/meet
  handler: internal/handler
  service: internal/service

swagger:
  path: doc
  code: code
  msg: msg
  data: data`
