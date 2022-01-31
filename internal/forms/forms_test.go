package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if !form.Valid() {
		t.Error("got invalid when should have been valid")
	}
}
func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b")
	if form.Valid() {
		t.Error("should failed, because we are requesting required values")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")

	form = New(postedData)
	form.Required("a", "b")
	if !form.Valid() {
		t.Error("it should be empty, cause we request with necessary params")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Has("a")
	if form.Valid() {
		t.Error("has should be invalid")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)
	form.Has("a")
	if !form.Valid() {
		t.Error("it should be empty, we add Has")
	}
}

func TestForm_MinLength(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("a", "a")
	form := New(postedData)
	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("minleng is nok")
	}

	postedData.Add("b", "aaaa")
	form = New(postedData)

	form.MinLength("b", 3)

	if !form.Valid() {
		t.Error("minleng should be ok")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}
	postedData.Add("a", "a")
	form := New(r.PostForm)
	form.IsEmail("a")
	if form.Valid() {
		t.Error("is not an email")
	}
	postedData.Add("a", "aa@aa.com")
	form = New(postedData)
	if !form.Valid() {
		t.Error("is an email")
	}
}

func TestErrors_Get(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("a", "aaa")
	form := New(postedData)

	form.MinLength("a", 3)
	if form.Errors.Get("a") != "" {
		t.Error("we should not have na error")
	}

	form.IsEmail("a")
	if form.Errors.Get("a") == "Not a valid Email" {
		t.Error("first error should be the email error")
	}
}
