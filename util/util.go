package util

import (
	"net/http"
	"log"
	"html/template"
)

type ErrorHandler struct {
	VersionInfo string
	Templates *template.Template
}

//HtmlErrorPage: write the appropriate html error page on an http error code
func (eh ErrorHandler) HtmlErrorPage(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.Printf("Writing an error message as html")
	log.Printf(err.Error())
	w.WriteHeader(code)
	d := map[string]interface{}{"Version": eh.VersionInfo,
		"StatusText": http.StatusText(code),
		"Msg": err.Error(),
		//"CASUser": map[bool]string{true: cas.Username(r), false: ""}[cas.IsAuthenticated(r)]
	}
	err = eh.Templates.ExecuteTemplate(w, "500.html.tmpl", d)
	if err != nil {
		log.Fatal(err)
	}
}

