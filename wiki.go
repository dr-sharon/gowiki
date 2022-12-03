package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)


// Data Structures
// The Page struct describes how page data will be stored in memory.
type Page struct {
	Title string 
	Body []byte
}


// Persistent storage method
// This method's signature reads: "This is a method named save that takes as its receiver p,
// a pointer to Page . It takes no parameters, and returns a value of type error." 
// This method will save the Page's Body to a text file.
// For simplicity, we will use the Title as the file name. 
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename,p.Body, 0600)
}


// The function loadPage constructs the file name from the title parameter,
// reads the file's contents into a new variable body, and returns a pointer to a Page literal
// constructed with the proper title and body values. 
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body},nil
}



func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err  := getTitle(w,r)
	if err != nil {
		return
	}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request){
	title, err  := getTitle(w,r)
	if err != nil {
		return
	}
	p ,err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// The function saveHandler will handle the submission of forms located on the edit pages.
// The page title (provided in the URL) and the form's only field, Body, are stored in a new Page. 
// The save() method is then called to write the data to a file, and the client is redirected to the /view/ page. 
// The value returned by FormValue is of type string. We must convert that value to []byte before it will fit into the Page struct.
//  We use []byte(body) to perform the conversion. 
func saveHandler(w http.ResponseWriter, r *http.Request){
	title, err  := getTitle(w,r)
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	// Any errors that occur during p.save() will be reported to the user. 
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// If we were to add more templates to our program, we would add their names to the ParseFiles call's arguments. 
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	//The http.Error function sends a specified HTTP response code 
	// (in this case "Internal Server Error") and error message. 
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // the title is the second subexpression.
}


func handler(w http.ResponseWriter, r *http.Request ) {
	fmt.Fprintf(w, "Hi there, I Love this go %s!", r.URL.Path[1:])
}


func main() {
	//http.HandleFunc("/", handler) 
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080",nil))
	// p1 := &Page{Title: "TestPage", Body: []byte("Hey this is a test page")}
	// p1.save()
	// p2, _ := loadPage("TestPage")
	// fmt.Println(string(p2.Body))


}
