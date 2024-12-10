package handlers

import (
	"github.com/gouef/router"
	"html/template"
	"log"
)

type TemplateHandler struct {
	Router *router.Router
}

// Typ pro funkci urlFor
type UrlForFunc func(name string, params ...interface{}) string

// Vytvoření a inicializace šablon
func (t *TemplateHandler) Initialize() {
	n := t.Router.GetNativeRouter()
	n.SetFuncMap(template.FuncMap{
		"snippet":    snippetStart,
		"snippetEnd": snippetEnd,
		"endSnippet": snippetEnd,
		"link": UrlForFunc(func(name string, params ...interface{}) string {
			// Převod parametrů do mapy
			paramMap := make(map[string]interface{})
			if len(params) > 0 {
				for i := 0; i < len(params); i += 2 {
					key := params[i].(string)
					value := params[i+1]
					paramMap[key] = value
				}
			}

			url, err := t.Router.GenerateUrlByName(name, paramMap)
			if err != nil {
				log.Println("Error generating URL for", name, ":", err)
				return "/error"
			}
			return url
		}),
	})
}
