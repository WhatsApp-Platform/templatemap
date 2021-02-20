/*
Package templatemap implements a map of names to html templates
and utility functions to load them from directories.
*/
package templatemap

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Map represents a map from template names to
// template.Template instances
type Map map[string]*template.Template

// Parser contains the configurations for parsing a
// template map
type Parser struct {
	// The template.FuncMap that will be added to
	// each template.Template instance in the map
	FuncMap template.FuncMap

	// A list of options that will be passed to
	// the template.Template.Option function on
	// each template.Template instance on the map
	Options []string
}

// Parser.ParseFS loads the templates from a given
// fs.FS recursively.
// Each template is added to the map as it's relative path
// from the directory being loaded.
// Each directory may also contain a "super.tmpl" file,
// wich will be added as an associated template to the
// directories and its childrens templates.
func (p *Parser) ParseFS(f fs.FS) (Map, error) {
	m := make(Map)
	err := p.parseFS(f, m, nil, ".", "")
	return m, err
}

// Parser.ParseDir does the same as Parser.ParseFS, but uses an
// operating system directory instead.
func (p *Parser) ParseDir(path string) (Map, error) {
	return p.ParseFS(os.DirFS(path))
}

func (p *Parser) parseFS(f fs.FS, m Map, super *template.Template, path, rel string) error {
	basePath := filepath.Join(path, "super.tmpl")
	baseName := rel + "super.tmpl"

	super, err := p.parseTemplate(f, super, basePath, baseName)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	entries, err := fs.ReadDir(f, path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		tmplPath := filepath.Join(path, entry.Name())
		tmplName := rel + entry.Name()

		if entry.IsDir() {
			err = p.parseFS(f, m, super, tmplPath, tmplName+"/")
			if err != nil {
				return err
			}
			continue
		}

		if filepath.Ext(entry.Name()) != ".tmpl" || entry.Name() == "super.tmpl" {
			continue
		}

		tmpl, err := p.parseTemplate(f, super, tmplPath, tmplName)
		if err != nil {
			return err
		}
		m[tmplName] = tmpl
	}
	return nil
}

func (p *Parser) parseTemplate(f fs.FS, super *template.Template, path, name string) (*template.Template, error) {
	file, err := f.Open(path)
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	_, err = io.Copy(&b, file)
	if err != nil {
		return nil, err
	}

	var tmpl *template.Template
	if super == nil {
		tmpl = template.New(name)
	} else {
		superC, err := super.Clone()
		if err != nil {
			// Impossible error?
			return nil, err
		}
		tmpl = superC.New(name)
	}

	if p.Options != nil {
		tmpl.Option(p.Options...)
	}
	if p.FuncMap != nil {
		tmpl.Funcs(p.FuncMap)
	}
	return tmpl.Parse(b.String())
}

// ParseFS loads the templates from a given fs.FS recursively.
// Each template is added to the map as it's relative path
// from the directory being loaded.
// Each directory may also contain a "_base.tmpl" file,
// wich will be added as an associated template to the
// directories and its childrens templates.
func ParseFS(f fs.FS) (Map, error) {
	var p Parser
	return p.ParseFS(f)
}

// ParseDir does the same as ParseFS, but uses an operating
// system directory instead.
func ParseDir(path string) (Map, error) {
	var p Parser
	return p.ParseDir(path)
}
