package browser

import (
	"github.com/pkg/browser"
)

// To open the url in browser.
func OpenBrowser(url string) error {
	return browser.OpenURL(url)
}
