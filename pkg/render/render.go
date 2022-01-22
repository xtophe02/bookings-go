package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/xtophe02/bookings-go/pkg/config"
	"github.com/xtophe02/bookings-go/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

//SETS CONFIG PKG HERE AVAILABLE
func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData)*models.TemplateData{
return td
}

func RenderTemplate(rw http.ResponseWriter, tmpl string, td *models.TemplateData){
	
	var tc map[string]*template.Template

	//DEV MODE: RUNS CREATE TEMPLATE ALL THE TIME, PROD MODE: USES TEMPLATES ONCE
	if app.UseCache{
		//GET ALL TEXT-TEMPLATES FROM APP CONFIG
		tc = app.TemplateCache
	} else {
			tc, _ = CreateTemplateCache()

	}
	

	//TEMPLATE WITH CODE INSERTED
	t, ok := tc[tmpl]
	if !ok{
		log.Fatal("no template")
	}

	//WE NEED TO READ THE TEMPLATE FROM DISK
	buf := new(bytes.Buffer)
	//HOLD INFO FROM MEMORY (CURRENT TEMPLATE) INTO A BUFFER OF BYTES
	t.Execute(buf, AddDefaultData(td))

	_, err := buf.WriteTo(rw)
	if err != nil{
		fmt.Println("error writing template to browser", err)
	}
}

//CREATES TEMPLATE CACHE AS A MAP
func CreateTemplateCache() (map[string]*template.Template,error){
	//MAP OF STRINGS AND VALUE OF TEXT-TEMPLATES
	myCache :=map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache,err
	}
	for _, page := range pages {
		//FETCH NAME OF FILE
		name := filepath.Base(page)
		//FETCH TEXT-TEMPLATE OF THE NAMED FILE
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache,err
		}

		//CHECKS IF A LAYOUT EXISTS
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache,err
		}
		if len(matches) > 0{
			//INJECTS THE LAYOUT ON TEXT-TEMPLATE
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache,err
			}
		}
		myCache[name] = ts
	}
	//RETURNS MAP OF NAMED FILED - TEXT-TEMPLATE 
	return myCache,nil
}