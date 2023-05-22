package luthor

import (
	"strings"

	goose "github.com/advancedlogic/GoOse"
	"golang.org/x/net/html"
)

/*
Parser is the top level object which directs the parsing and extraction of
information from a string of HTML.
*/
type Parser struct {
	html string
}

/*
NewParser creates a new Parser.
*/
func NewParser(html string) *Parser {
	return &Parser{html}
}

/*
GetElements returns a slice of strings containing all the instances
found of the element type that is passed in.
*/
func (p *Parser) GetElements(element string) ([]*html.Node, error) {
	elements := []*html.Node{}

	doc, err := html.Parse(strings.NewReader(p.html))
	if err != nil {
		return nil, err
	}

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == element {
			elements = append(elements, node)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	traverse(doc)

	return elements, nil
}

/*
GetArticle returns the article content
*/
func (p *Parser) GetArticle(result *string) error {
	g := goose.New()
	var extract string
	article, _ := g.ExtractFromRawHTML(p.html, extract)
	println("title", article.Title)
	println("description", article.MetaDescription)
	println("keywords", article.MetaKeywords)
	println("content", article.CleanedText)
	println("url", article.FinalURL)
	println("top image", article.TopImage)

	extract = article.Doc.Text()
	result = &extract
	return nil
}
