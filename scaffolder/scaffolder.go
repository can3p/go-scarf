package scaffolder

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var DefaultFuncMap template.FuncMap = map[string]any{
	"template_code": func(s string) template.HTML {
		return template.HTML(s)
	},
}

type ScaffoldData map[string]any

type OutputProcessor func(originalPath string, finalPath string, isDir bool, content []byte) error

func DebugProcessor(originalPath string, finalPath string, isDir bool, content []byte) error {
	fmt.Printf("path=%q, isDir=%v\n", originalPath, isDir)

	if !isDir {
		fmt.Printf("path=%q, output:\n\n%s\n\n", finalPath, string(content))
	}

	return nil
}

func FSProcessor(outputFolder string) OutputProcessor {
	return func(originalPath string, finalPath string, isDir bool, content []byte) error {
		writePath := outputFolder + "/" + finalPath

		folderToCreate := writePath

		if !isDir {
			folderToCreate = filepath.Dir(writePath)
		}

		// check file exitance there
		_ = os.MkdirAll(folderToCreate, 0777)

		if !isDir {
			return os.WriteFile(writePath, content, 0666)
		}

		return nil
	}
}

type scaffolder struct {
	funcMap         template.FuncMap
	outputProcessor OutputProcessor
}

func New() *scaffolder {
	return &scaffolder{
		funcMap:         DefaultFuncMap,
		outputProcessor: DebugProcessor,
	}
}

func (s *scaffolder) WithFuncMap(f template.FuncMap) *scaffolder {
	return &scaffolder{
		funcMap:         f,
		outputProcessor: s.outputProcessor,
	}
}

func (s *scaffolder) WithProcessor(p OutputProcessor) *scaffolder {
	return &scaffolder{
		funcMap:         s.funcMap,
		outputProcessor: p,
	}
}

func (s *scaffolder) Scaffold(tpl fs.FS, data ScaffoldData) error {
	return fs.WalkDir(tpl, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." || !d.IsDir() && !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		newName := strings.Replace(path, ".tmpl", "", -1)
		var buf bytes.Buffer

		if !d.IsDir() {
			content, err := fs.ReadFile(tpl, path)

			if err != nil {
				return err
			}

			tmpl := template.New(path).Funcs(DefaultFuncMap)
			tmpl, err = tmpl.Parse(string(content))

			if err != nil {
				return err
			}

			err = tmpl.ExecuteTemplate(&buf, path, data)

			if err != nil {
				return err
			}
		}

		return s.outputProcessor(path, newName, d.IsDir(), buf.Bytes())
	})
}
