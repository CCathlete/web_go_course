package olderversions

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>Contact page</h1><p>To get in touch, email me at <a"+
		"href=\"mailto:ccathlete01@gmail.com\">ccathlete01@gmail.com\n</a>.</p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	// Creating a data bucket for the content.
	var qaYaml struct {
		Questions string `yaml:"Content"`
	}

	// Opening the file for reading only.
	file, err := os.Open("QA.yaml")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while opening internal QA file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Spilling the data into the bucket in small chunks.
	if err := yaml.NewDecoder(file).Decode(&qaYaml); err != nil && err != io.EOF {
		http.Error(w, fmt.Sprintf("Error while opening internal QA file: %v", err), http.StatusInternalServerError)
		return
	}

	formattedContent := strings.ReplaceAll(html.EscapeString(qaYaml.Questions), "\n", "<br>")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>QA page</h1><p>%s</p>", formattedContent)
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found.", http.StatusNotFound)

	}

	fmt.Fprintf(w, "\nThe current path is %s", r.URL.Path)
}

type Router struct{} // A simple implementation to the http.Handler
// interface. It just needs to have the ServeHTTP method/

func (Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathHandler(w, r)
}

func main() {
	var myRouter Router
	// http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/contact", contactHandler)
	// http.HandleFunc("/", pathHandler)
	fmt.Println("Starting the server on: 3000...")
	http.ListenAndServe(":3000", myRouter)
}
