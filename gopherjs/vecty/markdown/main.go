package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func main() {
	vecty.SetTitle("Markdown Demo")
	vecty.RenderBody(&PageView{
		Input: `# Markdown Example

This is a live editor, try editing the Markdown on the right of the page.
`,
	})
}

type PageView struct {
	vecty.Core
	Input string
}

func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Div(
			vecty.Markup(
				vecty.Style("float", "right"),
			),
			elem.TextArea(
				vecty.Markup(
					vecty.Style("font-family", "monospace"),
					vecty.Property("rows", 14),
					vecty.Property("cols", 70),

					event.Input(func(e *vecty.Event) {
						p.Input = e.Target.Get("value").String()
						vecty.Rerender(p)
					}),
				),
				vecty.Text(p.Input),
			),
		),
		&Markdown{Input: p.Input},
	)
}

type Markdown struct {
	vecty.Core
	Input string `vecty:"prop"`
}

func (m *Markdown) Render() vecty.ComponentOrHTML {
	unsafeHTML := blackfriday.MarkdownCommon([]byte(m.Input))

	safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

	return elem.Div(
		vecty.Markup(
			vecty.UnsafeHTML(string(safeHTML)),
		),
	)
}
