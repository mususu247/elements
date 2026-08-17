# elements

## 階層型 JSON module


## インポートのテスト
```Bush
go get github.com/mususu247/elements@latest
go mod tidy
```

```Go
package main

import "github.com/mususu247/elements/elements"

func main() {
	var doc elements.Dom
	doc.Init("\t")
	root := doc.Root()
	docType := doc.CreateComment("<!DOCTYPE html>")
	root.AppendChild(docType)

	html := doc.CreateElement("html")
	html.SetAttribute("lang", "ja")
	root.AppendChild(html)

	head := doc.CreateElement("head")
	body := doc.CreateElement("body")
	html.AppendChild(head)
	html.AppendChild(body)

	doc.SaveAs("test.json")
	root.Export("test.html", ".html")
	doc.Close()
}
```

```html:test.html
<!DOCTYPE html>
<html lang="ja">
	<head></head>
	<body></body>
</html>
```

