package repo

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/longkai/xiaolongtongxue.com/github"
	"github.com/longkai/xiaolongtongxue.com/helper"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var re = regexp.MustCompile(`<article[\s\S]*>[\s\S]*</article>`)

func extractArticle(in io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	res := re.Find(b)
	if res == nil {
		return nil, fmt.Errorf("<article> not found in HTML")
	}
	return res, nil
}

// Renderer render a document into HTML.
type Renderer interface {
	Render(doc Doc) (template.HTML, error)
}

type newRender struct {
	wrap Renderer
}

func (r *newRender) Render(doc Doc) (template.HTML, error) {
	//if !strings.HasSuffix(doc.Path, ".md") {
	// TODO: use it for image cdn...
	if true {
		return r.wrap.Render(doc)
	}
	b, err := github.Markdown(strings.NewReader(doc.rawBody))
	var buf bytes.Buffer
	buf.WriteString(`<article class="markdown-body entry-content" itemprop="text">`)
	buf.Write(b)
	buf.WriteString(`</article>`)
	return template.HTML(buf.String()), err
}

// GithubRenderer render a mark up document using the Github favored style.
type GithubRenderer struct {
	User, Repo     string
	Dir            Dir
	StripTitle     bool
	URLTransformer func(src string) string
}

// NewRenderer default Github renderer, the URL in the HTML will remain unchanged.
func NewRenderer(user, repo string, dir Dir) Renderer {
	return &GithubRenderer{
		User:           user,
		Repo:           repo,
		Dir:            dir,
		StripTitle:     true,
		// TODO: hard code right now, should redesign next iteration.
		URLTransformer: func(src string) string { return "//cdn.jsdelivr.net/gh/longkai/essays" + src },
	}
}

// Render a file in a repository of a Github user.
func (r *GithubRenderer) Render(d Doc) (template.HTML, error) {
	const branch = "master" // May be support other branches other than master?
	f := r.Dir.Rel(d.Path)
	url := fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", r.User, r.Repo, branch, f[:strings.LastIndex(f, filepath.Ext(f))]+".md")

	resp, err := http.Get(url)
	if resp.StatusCode == http.StatusNotFound {
		url = fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", r.User, r.Repo, branch, f)
		resp, err = http.Get(url)
	}
	dummy := template.HTML("")
	if err != nil {
		return dummy, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dummy, fmt.Errorf("http status %d: %s", resp.StatusCode, url)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dummy, err
	}

	b, err = extractArticle(bytes.NewReader(b))
	if err != nil {
		return dummy, err
	}

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return dummy, err
	}
	// The stdlib encloses `html>body` for us.
	article := doc.FirstChild.LastChild.FirstChild
	return r.render(article)
}

func (r *GithubRenderer) render(article *html.Node) (template.HTML, error) {
	// Return -1 not found, otherwise the index.
	getAttr := func(n *html.Node, query func(a html.Attribute) bool) int {
		for i, v := range n.Attr {
			if query(v) {
				return i
			}
		}
		return -1
	}

	// Strip empty text element version.
	searchNextSibling := func(n *html.Node, query func(node *html.Node) bool) *html.Node {
		for c := n; c != nil; c = c.NextSibling {
			if t := c.DataAtom; t != atom.Atom(0) {
				if query(c) {
					return c
				}
				return nil
			}
		}
		return nil
	}

	if r.StripTitle {
		if h1 := searchNextSibling(article.FirstChild, func(node *html.Node) bool {
			if t := node.DataAtom; t == atom.H1 {
				return true
			}
			return false
		}); h1 != nil {
			article.RemoveChild(h1)
		}
	}

	queryYaml := func(node *html.Node) bool {
		if t := node.DataAtom; t == atom.Div {
			if getAttr(node, func(a html.Attribute) bool {
				return a.Key == "class" && strings.Contains(a.Val, "highlight-source-yaml")
			}) > -1 {
				return true
			}
		}
		return false
	}

	// Strip adjacent `EOF` and `yaml` block if any.
strip:
	for c := article.LastChild; c != nil; c = c.PrevSibling {
		if t := c.DataAtom; atom.H1 <= t && t <= atom.H6 {
			for test := c.FirstChild; test != nil; test = test.NextSibling {
				if test.Type == html.TextNode && strings.ToUpper(test.Data) == "EOF" {
					if yaml := searchNextSibling(c.NextSibling, queryYaml); yaml != nil {
						article.RemoveChild(c)
						article.RemoveChild(yaml)
						break strip
					}
				}
			}
		}
	}

	parseURL := func(url string) string {
		// Source file: `/{user}/{repo}/blob/master/path/to/file`
		s := fmt.Sprintf("/%s/%s/blob/master", r.User, r.Repo)
		if !strings.HasPrefix(url, s) {
			// Raw file: `/{user}/{repo}/raw/master/path/to/file`
			s = fmt.Sprintf("/%s/%s/raw/master", r.User, r.Repo)
			if !strings.HasPrefix(url, s) {
				// Patterns not match, so we cannot parse.
				return url
			}
		}
		return url[len(s):]
	}

	// Normalize `href` and `src` attribute`.
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			hrefIdx := getAttr(n, func(a html.Attribute) bool {
				return a.Key == "href" && strings.HasPrefix(a.Val, "/")
			})

			if img := n.FirstChild; img != nil && img.DataAtom == atom.Img {
				// The <a> followed by an <img />.
				if idx := getAttr(img, func(a html.Attribute) bool {
					// Only care relative URLs since Github will turn all the relative links into links starting with '/'.
					// If the link links to external site, github will cache it with its own CDN.
					return a.Key == "src" && strings.HasPrefix(a.Val, "/")
				}); idx > -1 {
					attr := img.Attr[idx]
					url := parseURL(attr.Val)
					img.Attr[idx] = html.Attribute{Namespace: attr.Namespace, Key: attr.Key, Val: r.URLTransformer(url)}

					// Apply back to <a> if any.
					if hrefIdx > -1 {
						href := n.Attr[hrefIdx]
						if i := strings.LastIndex(url, "."); i > -1 {
							// Test existence of a `img@full.ext`.
							urlFull := fmt.Sprintf("%s@full%s", url[:i], url[i:])
							if path := r.Dir.Abs(urlFull); helper.Exists(path) {
								url = urlFull
							}
						}
						n.Attr[hrefIdx] = html.Attribute{Namespace: href.Namespace, Key: href.Key, Val: r.URLTransformer(url)}
					}
				}
			} else {
				// Barely a link, no <img /> inside.
				if hrefIdx > -1 {
					href := n.Attr[hrefIdx]
					url := parseURL(href.Val)
					n.Attr[hrefIdx] = html.Attribute{Namespace: href.Namespace, Key: href.Key, Val: r.URLTransformer(url)}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(article)

	var buf bytes.Buffer
	if err := html.Render(&buf, article); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(buf.String()), nil
}
