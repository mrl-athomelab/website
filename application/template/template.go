package template

import (
	tpl "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrl-athomelab/website/application/logger"
)

var (
	funcMap = tpl.FuncMap{
		"add":         AddFunction,
		"as_html":     AsHTML,
		"format_time": FormatTime,
	}
)

func New(rootPath string) *tpl.Template {
	t := tpl.New("").Funcs(funcMap)
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = t.ParseFiles(path)
			if err != nil {
				logger.Warn("Template error %v", err)
			}
		}
		return err
	})

	return t
}
