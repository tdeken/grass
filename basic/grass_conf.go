package basic

type GrassConf struct {
	ModName string `yaml:"mod_name"` // mod 名称
	Proto   struct {
		Path     string `yaml:"path"`      // 存储目录
		FileType string `yaml:"file_type"` // 文件类型json，toml
	} `yaml:"proto"` // 协议源
	Analyze struct {
		Sources string `yaml:"sources"` // 参数资源存储目录
		Handler string `yaml:"handler"` // 处理方法存储目录
		Service string `yaml:"service"` // 依赖服务存储目录
	} `yaml:"analyze"` // 解析文件存储参数
	Swagger struct {
		Path string `yaml:"path"` // 存放目录
		Code string `yaml:"code"` // 返回的固定结构code
		Msg  string `yaml:"msg"`  // 返回的固定结构msg
		Data string `yaml:"data"` // 返回的固定结构data
	} `yaml:"swagger"` // 文档文件
}
