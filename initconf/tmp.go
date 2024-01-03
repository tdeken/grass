package initconf

const ConfTmp = `server:
  # listen prot
  port: 8888

conf_file:
  # options: path, env
  #   "path": The configuration file path is provided by the framework
  #   "env"： The configuration file path comes from a custom environment variable
  load_type: "path"

  # depend on load_type
  #   "path", options: local dev test prod
  #   "env", your environment variable name
  env: "local"

  # disable or enable hot load
  hot: "false"

  # file type
  type: "yaml"
`

const GrassTmp = `mod_name: {{ .ModName }}

source: proto

analyze:
  params: api/http
  entrance: internal/api
  service: 

swagger:
  path: api/swagger
  code: code
  msg: msg
  data: data`
