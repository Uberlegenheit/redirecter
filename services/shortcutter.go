package services

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
)

const (
	httpScheme  = "http"
	httpsScheme = "https"
)

type Data struct {
	TagCode     string
	URLRedirect string
}

func (s *Service) GetShortCutInfo(linkCode string) (buffer bytes.Buffer, err error) {
	link, found, err := s.dao.GetLinkByCode(linkCode)
	if err != nil {
		return s.getDefaultRedirectTemplate(), err
	}

	if !found {
		return s.getDefaultRedirectTemplate(), fmt.Errorf("link not found, probably fraud")
	}

	redirectURL := url.URL{
		Scheme: httpsScheme,
		Host:   s.conf.RedirectHost,
		Path:   "",
	}

	values := redirectURL.Query()
	values.Add("utm_source", link.Code)
	redirectURL.RawQuery = values.Encode()

	dt := Data{TagCode: link.Code, URLRedirect: redirectURL.String()}

	err = s.templates.ExecuteTemplate(&buffer, "redirect_body.html", dt)
	if err != nil {
		return s.getDefaultRedirectTemplate(), err
	}

	return buffer, nil
}

func (s *Service) getDefaultRedirectTemplate() (buffer bytes.Buffer) {

	err := s.templates.ExecuteTemplate(&buffer, "default_redirect_body.html", Data{URLRedirect: s.conf.DefaultRedirectURL})
	if err != nil {
		log.Errorf("getDefaultRedirectTemplate error: %s", err.Error())
	}

	return buffer
}
