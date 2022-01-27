package render

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"text/template"

	"github.com/justinas/nosurf"
	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

//SETS CONFIG PKG HERE AVAILABLE
func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request)*models.TemplateData{
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

func RenderTemplate(rw http.ResponseWriter, r *http.Request ,tmpl string, td *models.TemplateData){
	
	var tc map[string]*template.Template

	//DEV MODE: RUNS CREATE TEMPLATE ALL THE TIME, PROD MODE: USES TEMPLATES ONCE
	if app.InProduction{
		//GET ALL TEXT-TEMPLATES FROM APP CONFIG
		
		tc = app.TemplateCache
	} else {
			tc, _ = CreateTemplateCache()

	}
	

	//TEMPLATE WITH CODE INSERTED
	t, ok := tc[tmpl]
	if !ok{
		fmt.Println(app.InProduction)
		log.Fatal("no template")
	}

	//WE NEED TO READ THE TEMPLATE FROM DISK
	buf := new(bytes.Buffer)
	//HOLD INFO FROM MEMORY (CURRENT TEMPLATE) INTO A BUFFER OF BYTES
	t.Execute(buf, AddDefaultData(td,r))

	_, err := buf.WriteTo(rw)
	if err != nil{
		fmt.Println("error writing template to browser", err)
	}
}
var pages []string
func walk(s string, d fs.DirEntry, err error)  error {
	// fmt.Println(s)
	// fmt.Println(d)
	if err != nil {
		 return err
	}
	if ! d.IsDir() {
		// fmt.Println(s)
		pages = append(pages, s)
		
	}
	return nil
}


//CREATES TEMPLATE CACHE AS A MAP
func CreateTemplateCache() (map[string]*template.Template,error){
	//MAP OF STRINGS AND VALUE OF TEXT-TEMPLATES
	myCache :=map[string]*template.Template{}

	filepath.WalkDir(pathToTemplates,walk)

	// pages, err := filepath.Glob("./templates/*.page.tmpl")
	// if err != nil {
	// 	return myCache,err
	// }
// fmt.Println(pages)
	for _, page := range pages {
		// fmt.Println(page)
		//FETCH NAME OF FILE
		name := filepath.Base(page)
		
		//FETCH TEXT-TEMPLATE OF THE NAMED FILE
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		
		if err != nil {
			return myCache,err
		}

		//CHECKS IF A LAYOUT EXISTS
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
		if err != nil {
			return myCache,err
		}
		if len(matches) > 0{
			//INJECTS THE LAYOUT ON TEXT-TEMPLATE
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
			if err != nil {
				return myCache,err
			}
		}
		myCache[name] = ts
	}
	//RETURNS MAP OF NAMED FILED - TEXT-TEMPLATE 
	return myCache,nil
}