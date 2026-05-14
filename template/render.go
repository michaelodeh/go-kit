package template

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type TemplateRender struct {
	globalFiles  []string
	basePath     string
	extension    string
	systemParams map[string]interface{}
	fs           embed.FS
}

func NewRenderer(fs embed.FS) (*TemplateRender, error) {

	return &TemplateRender{
		fs: fs,
	}, nil
}

func (r *TemplateRender) AddTemplateFiles(patterns ...string) (*template.Template, error) {
	patterns = append(r.globalFiles, patterns...)
	for i, p := range patterns {
		if !filepath.IsAbs(p) {
			p = filepath.Join(r.basePath, p)
		}

		patterns[i] = p
	}

	tmpl, err := template.ParseFS(r.fs, patterns...)

	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		return nil, err
	}
	for _, t := range tmpl.Templates() {
		fmt.Println(t.Name())
	}
	return tmpl, nil
}

func (r *TemplateRender) AddLGlobalFiles(patterns ...string) {
	r.globalFiles = append(r.globalFiles, patterns...)
}

func (r *TemplateRender) SetTemplatePath(path string) {
	r.basePath = path
}

func (r *TemplateRender) Render(w io.Writer, tmpl *template.Template, name string, data *map[string]interface{}) error {
	if data == nil {
		data = &map[string]interface{}{}
	}
	for k, v := range r.systemParams {
		if _, ok := (*data)[k]; !ok {
			(*data)[k] = v
		}
	}
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (r *TemplateRender) RenderWithFiles(w http.ResponseWriter, name string, data *map[string]interface{}, patterns ...string) error {
	tmpl, err := r.AddTemplateFiles(patterns...)
	if err != nil {
		return err
	}
	return r.Render(w, tmpl, name, data)
}
