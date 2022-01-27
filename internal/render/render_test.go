package render

import (
	"net/http"
	"testing"

	"github.com/xtophe02/bookings-go/internal/models"
)

func TestAddDefaultData(t *testing.T){
	var td models.TemplateData

	r,err := getSession()
	if err != nil{
		t.Error(err)
	}
	session.Put(r.Context(),"flash","123")
	res := AddDefaultData(&td,r)

	if res.Flash != "123" {
		t.Error("Context not set on Session")
	}
}

func getSession()(*http.Request, error){
	r,err := http.NewRequest("GET","/some-url",nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx,r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}