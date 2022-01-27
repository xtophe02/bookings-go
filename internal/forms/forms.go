package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct{
	url.Values
	Errors errors
}

//IF NO ERRRORS RETURNS TRUE
func (f *Form) Valid() bool{
	return len(f.Errors) == 0
}

//NEW INITIALIZES A FORM STRUCT
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required (fields ...string){
	for _, field := range fields {
		
		v := f.Get(field)
		if strings.TrimSpace(v) == ""{
		f.Errors.Add(field, "This field cannot be blank")
	}}
}

//CHECK IF FIELD IS NOT EMPTY
func (f *Form) Has(field string, r *http.Request) bool{
	x:= r.Form.Get(field)
	if x == ""{
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}
func (f *Form) MinLength(field string, length int ,r *http.Request) bool{
	x:= r.Form.Get(field)
	if len(x) < length{
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long",length))
		return false
	}
	return true
}
func (f *Form) IsEmail(field string){
	if !govalidator.IsEmail(f.Get(field)){
		f.Errors.Add(field,"Not a valid email")
	}
}