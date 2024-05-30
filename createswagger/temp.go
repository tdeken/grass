package createswagger

import "html/template"

type SwaggerDocTemp struct {
	ModuleName string
	Content    template.HTML
}

var swaggerDocTemp = `package {{ .ModuleName }}

{{ .Content }}
func docDesc() {

}`

type SwaggerFileTemp struct {
	ModuleName string
	Group      string
}

var swaggerFileTemp = `package {{ .ModuleName }}

type {{ .Group }} struct {

}`

type SwaggerTemp struct {
	Name           string
	GroupDesc      string
	Desc           string
	SecurityTitle  string
	ReqContentType string
	ResContentType string
	Body           string
	Req            string
	ResFormat      string
	Route          string
	Method         string
	Group          string
	Messages       []template.HTML
	Annotation     []template.HTML
}

var swaggerTemp = `
{{ range $value := .Messages }}
{{ $value }}
{{ end }}

// {{ .Name }}
// @Tags {{ .GroupDesc }}
// @Summary {{ .Desc }}
// {{if .SecurityTitle }}@Security {{ .SecurityTitle }} {{ end }}
// @accept {{ .ReqContentType }}
// @Produce {{ .ResContentType }}{{ range $value := .Annotation }}
{{ $value }}{{ end }}
// @Param data {{ .Body }} {{ .Req }} true "数据"
// @Success 200 {object} {{ .ResFormat }}
// @Router {{ .Route }} [{{ .Method }}]
func ({{ .Group }}) {{ .Name }}() {

}`
