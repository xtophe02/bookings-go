package main

import (
	"net/http"
	"os"
	"testing"
)

//FUNC CALLED BEFORE THE TESTS, THAN WILL PERFORM ALL TESTS AND WILL EXIT AT THE END
func TestMain(m *testing.M){
	os.Exit(m.Run())
}

type myHandler struct{

}

func (mh *myHandler) ServeHTTP (rw http.ResponseWriter, r *http.Request){
	
}