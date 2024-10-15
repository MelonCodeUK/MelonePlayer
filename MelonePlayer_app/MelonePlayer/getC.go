// get components lib

package MelonePlayer

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetComponents() {
	res, err := http.Get("https://github.com/MelonCodeUK/MelonePlayer/tree/main/MelonePlayer_app/another")
	if err != nil {
		Log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		Log.Fatal("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		Log.Fatal(err)
	}
	Log.Info(string(doc.Text()))
}
