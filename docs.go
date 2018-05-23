package opark_backend

import (
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"html/template"
	"net/http"
	"fmt"
)

const (
	Api               = "doc/api.markdown"
	PrivacyStatement  = "doc/privacy_statement.markdown"
	PushNotifications = "doc/push_notifications.markdown"
	UsageStatistics   = "doc/usage_statistics.markdown"
	WorkerApi         = "doc/worker_api.markdown"
)

type DocumentationPageData struct {
	Title string
	Body  template.HTML
}

func renderDocument(w http.ResponseWriter, which string) {
	var title string
	switch which {
	case Api:
		title = "Client-server API"
	case PrivacyStatement:
		title = "Privacy statement"
	case UsageStatistics:
		title = "Usage statistics"
	default:
		panic(fmt.Sprintf("Unhandled document (%s) in renderDocument", which))
	}
	content, err := ioutil.ReadFile(which)
	handleError(err, fmt.Sprintf("Unable to read file (%s) in renderDocument", which))

	body := template.HTML(github_flavored_markdown.Markdown(content))
	templateParams := DocumentationPageData{Title: title, Body: body}
	tmpl, err := template.ParseFiles("templates/document.html")
	handleError(err, "Unable to parse template in renderDocument")
	err = tmpl.Execute(w, templateParams)
	handleError(err, "Unable to excecute template in renderDocument")
}
