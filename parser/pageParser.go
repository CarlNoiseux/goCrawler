package parser

import (
	"fmt"
	"goCrawler/context"
	"goCrawler/storage/storageTypes"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	//"goCrawler/storage/storageTypes"
)

func ParsePageUrls(ctx context.Context, urls chan string) {
	for {
		url, ok := <-urls
		if !ok {
			log.Println("urls channel closed, terminating parser goroutine")
			return
		}

		parsePageUrl(ctx, url)
	}
}

func parsePageUrl(ctx context.Context, url string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered a panic in parsePageUrl", r)
			(*ctx.Storage).UpdateUrlStatus(url, storageTypes.Unchartable)
		}
	}()

	panic("whatever")

	log.Printf("parsing url: %s", url)

	resp, err := retrievePage(url)
	if err != nil {
		(*ctx.Storage).UpdateUrlsStatuses([]string{url}, storageTypes.Charted)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not close request body")
		}
	}()

	// Parse retrieved page
	nodes, err := parsePage(resp)
	if err != nil {
		return
	}

	// Extract links from html nodes
	urls := extractUrls(nodes, true, resp)

	// Filter links (remove already explored, duplicates, etc..)
	_, missing := (*ctx.Storage).UrlsExist(urls)
	log.Printf("Found %d new urls while crawling %s", len(url), url)

	// Add links to storage
	for _, foundUrl := range missing {
		(*ctx.Storage).AddUrl(foundUrl, storageTypes.Uncharted)
	}
	(*ctx.Storage).UpdateUrlStatus(url, storageTypes.Charted)
}

func retrievePage(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	return resp, err
}

func parsePage(resp *http.Response) (*html.Node, error) {
	nodes, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing as HTML: %v", err)
	}

	return nodes, err
}

func extractUrls(nodes *html.Node, breathFirst bool, resp *http.Response) []string {
	links := make(map[string]bool)
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}

				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}

				links[link.String()] = true
			}
		}
	}

	if breathFirst {
		breadthFirstTraversal(nodes, visitNode)
	} else {
		depthFirstTraversal(nodes, visitNode)
	}

	urls := make([]string, len(links))
	i := 0
	for k := range links {
		urls[i] = k
		i++
	}

	return urls
}

func breadthFirstTraversal(node *html.Node, nodeParser func(n *html.Node)) {
	// TODO: Not ideal, using a slice, underlying array will be resized many times, potentially as long as there are nodes in the tree. would need a more intelligent queue
	children := make([]*html.Node, 0)
	for c := node; c != nil; {
		nodeParser(c)
		if c.FirstChild != nil {
			children = append(children, c.FirstChild)
		}

		if c.NextSibling != nil {
			c = c.NextSibling
		} else {
			if len(children) > 0 {
				c = children[0]
				children = children[1:]
			} else {
				c = nil
			}
		}
	}
}

func depthFirstTraversal(node *html.Node, nodeParser func(n *html.Node)) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		depthFirstTraversal(c, nodeParser)
	}

	nodeParser(node)
}
