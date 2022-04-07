package handlers

import "net/http"

var StaticHandler = http.FileServer(http.Dir("template/"))
