# elements

## 階層型 JSON module


## インポートのテスト
```Bush
go get github.com/mususu247/elements@latest
go mod tidy
```

```Go
package main

import (
	"fmt"
	"strconv"
	"github.com/mususu247/elements/elements"
)

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
	html.AppendChild(head)

	body := doc.CreateElement("body")
	html.AppendChild(body)

	title := doc.CreateElement("title")
	title.InnerText("T.I.T.L.E.")
	head.AppendChild(title)

	table := doc.CreateElement("table")
	table.SetAttribute("border", 1)
	body.AppendChild(table)

	r0 := doc.CreateElement("tr")
	r0.InnerHTML("<th>A.</th><th>B.</th><th>C.</th>")
	r0.SetAttribute("id", "R0")
	table.AppendChild(r0)

	for i := range 5 {
		rx := doc.CreateElement("tr")
		text := fmt.Sprintf("<td>A%v</td><td>B%v</td><td>C%v</td>", i+1, i+2, i+3)
		rx.InnerHTML(text)
		rx.SetAttribute("id", "R"+strconv.Itoa(i+1))
		table.AppendChild(rx)
	}

	r3 := root.GetElementById("R3")
	r6 := r3.CloneNode(true)
	r6.SetAttribute("id", "R6")
	r3.InsertAdjacentElemnt("beforebegin", r6)

	doc.SaveAs("index.json")
	root.Export("index.html", ".html")
	doc.Close()
}
```

```html:index.html
<!DOCTYPE html>
<html lang="ja">

<head>
    <title>T.I.T.L.E.</title>
</head>

<body>
    <table border="1">
        <tr id="R0"><th>A.</th><th>B.</th><th>C.</th></tr>
        <tr id="R1"><td>A1</td><td>B2</td><td>C3</td></tr>
        <tr id="R2"><td>A2</td><td>B3</td><td>C4</td></tr>
        <tr id="R6"><td>A3</td><td>B4</td><td>C5</td></tr>
        <tr id="R3"><td>A3</td><td>B4</td><td>C5</td></tr>
        <tr id="R4"><td>A4</td><td>B5</td><td>C6</td></tr>
        <tr id="R5"><td>A5</td><td>B6</td><td>C7</td></tr>
    </table>
</body>

</html>
```

