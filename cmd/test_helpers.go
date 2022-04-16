package cmd

const (
	shortLength = 30
	longLength  = 256
)

const (
	xssHtmlOpen  = "<a href=\"javascript:alert('XSS1')\" onmouseover=\"alert('XSS2')\">"
	xssHtmlClose = "</a>"
)
