package parser

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type Parser struct {
	PackageName string
	ImportType  string
	Name        string
	ChannelType string
	template *template.Template
}

func parsePackage() (string, error) {
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	pks, err := packages.Load(nil, args...)
	if err != nil {
		return "", fmt.Errorf("load packeges failed: %w", err)
	}
	if len(pks) != 1 {
		return "", fmt.Errorf("once package not found, packages' count: %d", len(pks))
	}
	return pks[0].Name, nil
}

func (p *Parser) init() error {
	if p.ChannelType == "" {
		return fmt.Errorf("type must be set")
	}
	if p.Name == "" {
		p.Name = p.ChannelType
	}
	tpl, err := template.New(p.Name).Parse(tmpl)
	if err != nil {
		return err
	}
	p.template = tpl

	pkg, err := parsePackage()
	if err != nil {
		return err
	}
	p.PackageName = pkg
	return nil
}

func (p *Parser) Execute() (err error) {
	if err := p.init(); err != nil {
		return err
	}

	const fileNameSuffix = "_channel.go"
	file, err := os.Create(strings.ToLower(p.Name) + fileNameSuffix)
	if err != nil {
		return fmt.Errorf("create go-file failed: %w", err)
	}
	defer func(err *error) {
		if fErr := file.Close(); fErr != nil {
			*err = fErr
		}
	}(&err)

	if err := p.template.Execute(file, p); err != nil {
		return err
	}
	return nil
}

func New(name, channelType, importType string) *Parser {
	return &Parser{
		Name: name,
		ChannelType: channelType,
		ImportType: importType,
	}
}
