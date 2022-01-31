package render

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/xtophe02/bookings-go/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	res := AddDefaultData(&td, r)

	if res.Flash != "123" {
		t.Error("Context not set on Session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, _ := http.NewRequest("GET", "/some-url", nil)

	var ww myWriter
	err = RenderTemplate(&ww, r, "non-existing.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("render template that does not exist")
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}
func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
func TestWalk(t *testing.T) {
	pathToTemplates = "./../../templates"
	err := filepath.WalkDir(pathToTemplates, Walk)
	if err != nil {
		t.Error(err)
	}
}
