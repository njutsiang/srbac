package srbac

import "github.com/gin-gonic/gin"

var Routes = []Route{}

type Route struct{
	Method string
	Uri string
	Name string
}

func GET(path string, handler gin.HandlerFunc, names ...string) {
	name := ""
	if len(names) == 1 {
		name = names[0]
	}
	Engine.GET(path, handler)
	Routes = append(Routes, Route{
		Method: "GET",
		Uri: path,
		Name: name,
	})
}

func POST(path string, handler gin.HandlerFunc, names ...string) {
	name := ""
	if len(names) == 1 {
		name = names[0]
	}
	Engine.POST(path, handler)
	Routes = append(Routes, Route{
		Method: "POST",
		Uri: path,
		Name: name,
	})
}

func DELETE(path string, handler gin.HandlerFunc, names ...string) {
	name := ""
	if len(names) == 1 {
		name = names[0]
	}
	Engine.DELETE(path, handler)
	Routes = append(Routes, Route{
		Method: "DELETE",
		Uri: path,
		Name: name,
	})
}