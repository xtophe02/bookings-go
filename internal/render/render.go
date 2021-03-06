package render

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"text/template"

	"github.com/justinas/nosurf"
	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/models"
)

var functions = template.FuncMap{
	"humanDate":         HumanDate,
	"formatDate":        FormatDate,
	"iterate":           Iterate,
	"formatDateWeekDay": FormatDateWeekDay,
}

var app *config.AppConfig
var pathToTemplates = "./templates"

//SETS CONFIG PKG HERE AVAILABLE
func NewRenderer(a *config.AppConfig) {
	app = a
}
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}
func FormatDateWeekDay(t time.Time, d int) string {
	now := time.Date(t.Year(), t.Month(), d, 0, 0, 0, 0, time.UTC)

	return now.Weekday().String()[0:3]
}

func Iterate(count int) []int {
	var items []int
	for i := 1; i < count+1; i++ {
		items = append(items, i)
	}
	return items
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

func Template(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	//DEV MODE: RUNS CREATE TEMPLATE ALL THE TIME, PROD MODE: USES TEMPLATES ONCE
	if app.InProduction {
		//GET ALL TEXT-TEMPLATES FROM APP CONFIG

		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()

	}

	//TEMPLATE WITH CODE INSERTED
	t, ok := tc[tmpl]
	if !ok {
		fmt.Println(app.InProduction)
		return errors.New("no template")
	}

	//WE NEED TO READ THE TEMPLATE FROM DISK
	buf := new(bytes.Buffer)
	//HOLD INFO FROM MEMORY (CURRENT TEMPLATE) INTO A BUFFER OF BYTES
	t.Execute(buf, AddDefaultData(td, r))

	_, err := buf.WriteTo(rw)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}

var pages []string

func Walk(s string, d fs.DirEntry, err error) error {
	// fmt.Println(s)
	// fmt.Println(d)
	if err != nil {
		return err
	}
	if !d.IsDir() {
		// fmt.Println(s)
		pages = append(pages, s)

	}
	return nil
}

//CREATES TEMPLATE CACHE AS A MAP
func CreateTemplateCache() (map[string]*template.Template, error) {
	//MAP OF STRINGS AND VALUE OF TEXT-TEMPLATES
	myCache := map[string]*template.Template{}

	_ = filepath.WalkDir(pathToTemplates, Walk)

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
			return myCache, err
		}

		//CHECKS IF A LAYOUT EXISTS
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			//INJECTS THE LAYOUT ON TEXT-TEMPLATE
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	//RETURNS MAP OF NAMED FILED - TEXT-TEMPLATE
	return myCache, nil
}
