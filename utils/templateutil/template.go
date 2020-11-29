package templateutil

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TemplateDirsEnv is environment key for templates directory
const TemplateDirsEnv = "TEMPLATES_DIR"

// ParseTemplate parses a given template
func ParseTemplate(files ...string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// ReadFiles read files from a directory
func ReadFiles(dir string) ([]string, error) {
	var allFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if (strings.HasSuffix(filename, ".html") || strings.HasSuffix(filename, ".htm")) && !file.IsDir() {
			allFiles = append(allFiles, filepath.Join(dir, filename))
		}
	}
	return allFiles, nil
}

// EmailData contains data that is sent as part of email body
type EmailData struct {
	Names          string
	AccountID      string
	Link           string
	Token          string
	Label          string
	WebsiteURL     string
	AppName        string
	AppDescription string
	PrimaryColor   string
	SecondaryColor string
	TemplateName   string
	Reason         string
}

var (
	primaryColor   string
	secondaryColor string
	appName        string
	websiteURL     string
)

// SetDefaults sets defaults for email data
func SetDefaults(app, website, mainColor, secColor string) {
	appName = app
	websiteURL = website
	primaryColor = mainColor
	secondaryColor = secColor
}

// DefaultEmailData creates email data with defaults
func DefaultEmailData() *EmailData {
	return &EmailData{
		AppName:        appName,
		WebsiteURL:     websiteURL,
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
	}
}
