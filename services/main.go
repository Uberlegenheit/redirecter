package services

import (
	"everstake-affiliate/conf"
	"everstake-affiliate/dao"
	"html/template"
)

type Service struct {
	conf      conf.Config
	dao       dao.DbDAO
	templates *template.Template
}

func New(c conf.Config, d dao.DbDAO) (*Service, error) {

	parsedTemplates, err := template.ParseFiles(
		"resources/templates/redirect_body.html",
		"resources/templates/default_redirect_body.html",
	)
	if err != nil {
		return nil, err
	}

	s := &Service{
		conf:      c,
		templates: parsedTemplates,
		dao:       d,
	}

	return s, nil
}
