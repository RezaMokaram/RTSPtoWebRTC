package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ShowBody(req http.Request) {
	tmp, _ := ioutil.ReadAll(req.Body)
	fmt.Println("req body\n", string(tmp), "\n")
}