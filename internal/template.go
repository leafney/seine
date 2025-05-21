/**
 * @Project:     seine
 * @Date:        2025-05-21
 */

package internal

import (
	"embed"
	"path/filepath"
)

//go:embed template/*.tmpl
//go:embed template/**/*.tmpl
//go:embed template/**/**/*.tmpl
var TemplateFS embed.FS

// GetTemplateContent 获取模板文件内容
func GetTemplateContent(templateName string) (string, error) {
	content, err := TemplateFS.ReadFile(filepath.Join("../../template", templateName))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// //go:embed *
// //go:embed */*
// //go:embed */*/*
// var TemplateFS embed.FS

// // ReadTemplateFile 读取模板文件内容
// func ReadTemplateFile(templateName string) ([]byte, error) {
// 	return TemplateFS.ReadFile(templateName)
// }

// // GetTemplateFS 获取模板文件系统
// func GetTemplateFS() fs.FS {
// 	return TemplateFS
// }

// // ListTemplates 列出所有模板文件
// func ListTemplates() ([]string, error) {
// 	var templates []string
// 	err := fs.WalkDir(TemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !d.IsDir() && filepath.Ext(path) == ".tmpl" {
// 			templates = append(templates, path)
// 		}
// 		return nil
// 	})
// 	return templates, err
// }
