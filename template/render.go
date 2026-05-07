package template

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
)

type TemplateRender struct {
	tmpl         *template.Template
	systemParams map[string]interface{}
	fs           embed.FS
}

func NewRenderer(fs embed.FS, patterns ...string) (*TemplateRender, error) {
	tmpl, err := template.ParseFS(fs, patterns...)

	if err != nil {
		return nil, err
	}
	for _, t := range tmpl.Templates() {
		fmt.Println(t.Name())
	}
	return &TemplateRender{
		tmpl: tmpl,
		fs:   fs,
	}, nil
}

func (r *TemplateRender) Render(w io.Writer, name string, data *map[string]interface{}) error {
	if data == nil {
		data = &map[string]interface{}{}
	}
	for k, v := range r.systemParams {
		if _, ok := (*data)[k]; !ok {
			(*data)[k] = v
		}
	}
	err := r.tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
