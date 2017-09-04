package urioj

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Sample struct {
	Input  []string
	Output []string
}

func (p *Problem) Name() string {
	return strings.TrimSpace(p.doc.Find("div.header > h1").Text())
}

func (p *Problem) Description() []string {
	return extractContent(p.doc.Find("div.description"))
}

func (p *Problem) Input() []string {
	return extractContent(p.doc.Find("div.input"))
}

func (p *Problem) Output() []string {
	return extractContent(p.doc.Find("div.output"))
}

func (p *Problem) Samples() []Sample {
	samples := make([]Sample, 0, 5)
	table := p.doc.Find("tbody")
	for i := range table.Nodes {
		sample := table.Eq(i).Find("td")
		input := formatSample(sample.First().Text())
		output := formatSample(sample.Last().Text())
		samples = append(samples, Sample{input, output})
	}
	return samples
}

func removeRedundantChar(s string) string {
	reg := regexp.MustCompile(`\r?\n?\t`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	return reg.ReplaceAllString(s, " ")
}

func text(n *html.Node) []string {
	var text []string

	var str string
	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Data {
		case "br":
			if str != "\n" {
				text = append(text, strings.TrimSpace(str))
				str = "\n"
			}
		case "sup":
			if len(str) > 0 && str[len(str)-1:] == " " {
				str = str[:len(str)-1]
			}
			str += "^" + strings.TrimSpace(n.FirstChild.Data)
		default:
			if n.Type == html.TextNode {
				data := removeRedundantChar(n.Data)
				if len(strings.TrimSpace(data)) > 0 {
					str += data
				}
			} else if n.FirstChild != nil {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		}
	}

	f(n)
	return append(text, strings.TrimSpace(str))
}

func extractContent(s *goquery.Selection) []string {
	content := make([]string, 0, 10)

	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Data {
		case "p":
			content = append(content, text(n)...)
		case "pre":
			for _, t := range text(n) {
				content = append(content, "  "+t)
			}
		default:
			if n.FirstChild != nil {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return content
}

func formatSample(s string) []string {
	e := make([]string, 0, 5)
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for _, line := range lines {
		e = append(e, removeRedundantChar(strings.TrimSpace(line)))
	}
	return e
}
