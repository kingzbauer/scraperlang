package interpreter

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Selector is the main interface implemented when quering into a document e.g HTML document
type Selector interface {
	Accessor
}

// Noder is a single html node from which we can access node attributes
type Noder interface {
	GetAttribute(key string) string
}

// Selection implements the selector interface
type Selection struct {
	document *goquery.Document
}

// Node represents a single HTML node and implements the Noder interface
type Node struct {
	node *html.Node
}
