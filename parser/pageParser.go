package parser

import (
	"fmt"
	"goCrawler/context"
	"goCrawler/storage/storageTypes"
	"golang.org/x/net/html"
	"net/http"
	//"goCrawler/storage/storageTypes"
)

func ParsePageUrls(ctx context.Context, url string) {
	// Retrieve page
	resp, err := retrievePage(url)
	if err != nil {
		(*ctx.Storage).UpdateUrlsStatuses([]string{url}, storageTypes.Charted)
		return
	}
	defer resp.Body.Close()

	// Parse retrieved page
	nodes, err := parsePage(resp)
	if err != nil {
		return
	}

	// Extract links from html nodes
	urls := extractUrls(nodes, true, resp)

	// Filter links (remove already explored, duplicates, etc..)

	// Add links to storage
	for _, url := range urls {
		(*ctx.Storage).WriteUrl(url, storageTypes.Uncharted)
	}

	return
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
	var links []string
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

				links = append(links, link.String())
			}
		}
	}

	if breathFirst {
		breadthFirstTraversal(nodes, visitNode)
	} else {
		depthFirstTraversal(nodes, visitNode)
	}

	return links
}

func breadthFirstTraversal(node *html.Node, nodeParser func(n *html.Node)) {
	// TODO: Not ideal, using a slice, underlying array will be as long as there are nodes in the tree. would need a more intelligent queue
	children := make([]*html.Node, 0)
	for c := node; c != nil; {
		nodeParser(node)
		if c.FirstChild != nil {
			children = append(children, c.FirstChild)
		}

		if c.NextSibling != nil {
			c = c.NextSibling
		} else {
			c = children[0]
			children = children[1:]
		}
	}

}

func depthFirstTraversal(node *html.Node, nodeParser func(n *html.Node)) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		depthFirstTraversal(c, nodeParser)
	}

	nodeParser(node)
}
