package utils

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"os/exec"
	"strings"
)

// IsFileExist 判断文件文件夹是否存在
func IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}

	//我这里判断了如果是0也算不存在
	if !fileInfo.IsDir() && fileInfo.Size() == 0 {
		return false, nil
	}

	if err == nil {
		return true, nil
	}

	return false, err
}

// MkDirAll 连续创建目录
func MkDirAll(paths ...string) error {
	for _, path := range paths {
		if path == "" {
			continue
		}
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// NotExistCreateDir 不存在目录就创建目录
func NotExistCreateDir(path string) error {
	exist, err := IsFileExist(path)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	return os.Mkdir(path, 0777)
}

// NotExistCreateFile 不存在就创建文件
func NotExistCreateFile(filename, content string) error {
	exist, err := IsFileExist(filename)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	return CreateFile(filename, content)
}

// CamelString 蛇形转驼峰
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// GetClass 结构体解析
func GetClass(name, class string, dictionary map[string]string) string {
	if strings.Contains(class, "?") {
		spt := strings.Split(class, "?")

		spt[len(spt)-1] = "*" + name + spt[len(spt)-1]
		class = strings.Join(spt, "")
	}

	if strings.Contains(class, "~") {
		spt := strings.Split(class, "~")

		spt[len(spt)-1] = "*" + dictionary[spt[len(spt)-1]]
		class = strings.Join(spt, "")
	}

	if strings.Contains(class, "&") {
		spt := strings.Split(class, "&")

		spt[len(spt)-1] = name + spt[len(spt)-1]
		class = strings.Join(spt, "")
	}

	return class
}

// MidString 大小写之间用特定字符分割
func MidString(s string, sep byte) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, sep)
		}
		if d != sep {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

// SnHeader 蛇形转请求头参数
func SnHeader(s string) string {
	arr := strings.Split(s, "_")

	var str string
	for i, v := range arr {
		if i > 0 {
			str += "-"
		}
		str += strings.ToUpper(string(v[0])) + v[1:]
	}

	return str
}

// CreateTmp 创建模版数据
func CreateTmp(content any, parse string) (text string, err error) {
	t, err := template.New("tmp.tpl").Parse(parse)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, content)
	if err != nil {
		return
	}

	return buf.String(), nil
}

// CreateFile 创建文件
func CreateFile(filename, content string) (err error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)

	return
}

// AppendFile 追加内容到文件
func AppendFile(filename, content string) (err error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)

	return
}

// RunCommand 运行命令
func RunCommand(dir, name string, args ...string) (err error) {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		err = errors.New(out.String())
		return
	}

	return
}

// PkgAndStruct 得到包和边准包类型名称
func PkgAndStruct(path string) (pkg, str string) {
	pkg = path[strings.LastIndex(path, "/")+1:]

	str = CamelString(pkg)
	return
}
