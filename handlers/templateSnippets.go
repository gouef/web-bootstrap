package handlers

import "html/template"

var (
	inSnippet bool // Flag pro kontrolu, jestli je snippet otevřen
)

// Funkce pro generování začátku snippetu s unikátním ID
func snippetStart(name string) template.HTML {
	if inSnippet {
		// Pokud je snippet již otevřen, vrátíme chybu nebo ignorujeme volání
		panic("Error: snippet already opened.")
		return template.HTML("<p>Error: snippet already opened.</p>")
	}

	inSnippet = true
	return template.HTML("<div id=\"snippet-" + name + "\">")
}

// Funkce pro uzavření snippetu
func snippetEnd() template.HTML {
	if !inSnippet {
		panic("Error: snippet not started.")
		return template.HTML("<p>Error: snippet not started.</p>")
	}
	inSnippet = false
	return template.HTML("</div>")
}
