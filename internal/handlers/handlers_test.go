package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/xtophe02/bookings-go/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var theTests = []struct {
	name   string
	url    string
	method string
	// params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", //]postData{}
		http.StatusOK},
	{"about", "/about", "GET", //[]postData{}
		http.StatusOK},
	{"generals-quarters", "/rooms/general-quarters", "GET", //[]postData{},
		http.StatusOK},
	{"majors-suite", "/rooms/major-suite", "GET", //[]postData{},
		http.StatusOK},
	{"search-availability", "/availability", "GET", //[]postData{},
		http.StatusOK},
	{"contact", "/contact", "GET", //[]postData{},
		http.StatusOK},
	{"non-existent", "/green/eggs", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "GET", http.StatusOK},
	{"show reservations", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	// {"make-res", "/reservation", "GET", []postData{}, http.StatusOK},
	// {"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	// {"post-avalability", "/availability", "POST",[]postData{
	// 		{key: "start", value: "2020-01-01"},
	// 		{key: "end", value: "2020-01-02"},
	// 	},http.StatusOK},
	// {"post-reservation", "/reservation", "POST",
	// 	[]postData{
	// 		{key: "first_name", value: "John"},
	// 		{key: "last_name", value: "Doe"},
	// 		{key: "email", value: "john.doe@gmail.com"},
	// 	},		http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		//if e.method == "GET" {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

		//	}
		// else {
		// 	values := url.Values{}
		// 	for _, x := range e.params {
		// 		values.Add(x.key, x.value)
		// 	}
		// 	res, err := ts.Client().PostForm(ts.URL+e.url, values)
		// 	if err != nil {
		// 		t.Log(err)
		// 		t.Fatal(err)
		// 	}

		// 	if res.StatusCode != e.expectedStatusCode {
		// 		t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
		// 	}
		// }
	}
}

func TestRepository_Reservation(t *testing.T) {
	//necessaty date for testing
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//new request
	req, _ := http.NewRequest("GET", "/reservation", nil)
	//load session
	ctx := getCtx(req)
	//inject context with session
	req = req.WithContext(ctx)
	//simlutes an http
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//test case not existing room
	req, _ = http.NewRequest("GET", "/reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// reqBody := "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	postedDate := url.Values{}
	postedDate.Add("start_date", "2050-01-01")
	postedDate.Add("end_date", "2050-01-02")
	postedDate.Add("room_id", "1")

	//TEST PARSEFORM
	req, _ := http.NewRequest("POST", "/availability-json", strings.NewReader(postedDate.Encode()))
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal((rr.Body.Bytes()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func TestRepository_PostReservation(t *testing.T) {
	//create reservation put put on sesssion
	reservation := models.Reservation{
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now(),
		FirstName: "john",
		LastName:  "smith",
		Email:     "john.smith@gmail.com",
		Phone:     "1234567",
	}

	//test parse form
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=john")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//TEST PARSEFORM
	req, _ := http.NewRequest("POST", "/reservation", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//TEST SESSION RESERVATION
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//TEST Form invalid
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusSeeOther)
	}

	//FAIL RESERVATIONID INTO DB
	reservation.RoomID = 2
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//FAIL RRESTRICTION INTO DB
	reservation.RoomID = 1000
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//FINAL SEEOTHER
	reservation.RoomID = 1
	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusSeeOther)
	}

	// //new request
	// req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(reqBody))
	// //load session
	// ctx = getCtx(req)
	// //inject context with session
	// req = req.WithContext(ctx)
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// //simlutes an http
	// rr = httptest.NewRecorder()
	// handler = http.HandlerFunc(Repo.PostReservation)
	// // session.Put(ctx, "reservation", reservation)
	// handler.ServeHTTP(rr, req)
	// if rr.Code != http.StatusTemporaryRedirect {
	// 	t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusTemporaryRedirect)
	// }
	// rr = httptest.NewRecorder()
	// session.Put(ctx, "reservation", reservation)
	// handler.ServeHTTP(rr, req)
	// if rr.Code != http.StatusSeeOther {
	// 	t.Errorf("PostReservation handler returned wrong response code. got %d but wanted %d", rr.Code, http.StatusSeeOther)
	// }

}

func TestRepository_ChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//TEST WITHOUT SESSION
	// req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	// ctx = getCtx(req)
	// req = req.WithContext(ctx)
	// req.RequestURI = "/choose-room/1"
	// rr = httptest.NewRecorder()
	// handler.ServeHTTP(rr, req)
	// if rr.Code != http.StatusSeeOther {
	// 	t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	// }

}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{"valid-credentials", "me@here.ca", http.StatusSeeOther, "", "/"},
	{"invalid-credentials", "jack@nimble.com", http.StatusOK, `action="/user/login`, ""},
	{"invalid-DATA", "j", http.StatusOK, `action="/user/login"`, ""},
}

func TestLogin(t *testing.T) {
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		req = req.WithContext(getCtx(req))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d got %d", e.name, e.expectedStatusCode, rr.Code)
		}
		if e.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s got %s", e.name, e.expectedLocation, actualLocation.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s ", e.name, e.expectedHTML)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
