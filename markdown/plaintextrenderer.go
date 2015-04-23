package markdown

import (
	"bytes"
)

type PlaintextRenderer struct{}

func (options *PlaintextRenderer) GetFlags() int {
	return 0
}

func (options *PlaintextRenderer) TitleBlock(out *bytes.Buffer, text []byte) {
	text = bytes.TrimPrefix(text, []byte("% "))
	text = bytes.Replace(text, []byte("\n% "), []byte("\n"), -1)
	out.Write(text)
	out.WriteString("\n")
}

func (options *PlaintextRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()
	options.wrap(out)

	if !text() {
		out.Truncate(marker)
		return
	}
}

func (options *PlaintextRenderer) BlockHtml(out *bytes.Buffer, text []byte) {
	options.wrap(out)
	out.Write(text)
	out.WriteByte('\n')
}

func (options *PlaintextRenderer) HRule(out *bytes.Buffer) {
	options.wrap(out)
}

func (options *PlaintextRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	options.BlockCodeNormal(out, text, lang)
}

func (options *PlaintextRenderer) BlockCodeNormal(out *bytes.Buffer, text []byte, lang string) {
	options.wrap(out)
	out.Write(text)
}

func (options *PlaintextRenderer) BlockQuote(out *bytes.Buffer, text []byte) {
	options.wrap(out)
	out.Write(text)
}

func (options *PlaintextRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	options.wrap(out)
	out.Write(header)
	out.Write(body)
}

func (options *PlaintextRenderer) TableRow(out *bytes.Buffer, text []byte) {
	options.wrap(out)
	out.Write(text)
}

func (options *PlaintextRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	options.wrap(out)
	out.Write(text)
}

func (options *PlaintextRenderer) TableCell(out *bytes.Buffer, text []byte, align int) {
	options.wrap(out)
	out.Write(text)
}

func (options *PlaintextRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	options.HRule(out)
	options.List(out, text, 0)
}

func (options *PlaintextRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	out.Write(text)
}

func (options *PlaintextRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	options.wrap(out)

	if !text() {
		out.Truncate(marker)
		return
	}
}

func (options *PlaintextRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.Write(text)
	out.WriteString("\n")
}

func (options *PlaintextRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	options.wrap(out)

	if !text() {
		out.Truncate(marker)
		return
	}
}

func (options *PlaintextRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.Write(link)
}

func (options *PlaintextRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *PlaintextRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *PlaintextRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.Write(text)
}

func (options *PlaintextRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	return
}

func (options *PlaintextRenderer) LineBreak(out *bytes.Buffer) {
	return
}

func (options *PlaintextRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.Write(content)
	if !options.isRelativeLink(link) {
		out.WriteString(" ")
		out.Write(link)
	}
	return
}

func (options *PlaintextRenderer) RawHtmlTag(out *bytes.Buffer, text []byte) {
	return
}

func (options *PlaintextRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *PlaintextRenderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *PlaintextRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	return
}

func (options *PlaintextRenderer) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (options *PlaintextRenderer) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *PlaintextRenderer) Smartypants(out *bytes.Buffer, text []byte) {}

func (options *PlaintextRenderer) DocumentHeader(out *bytes.Buffer) {}

func (options *PlaintextRenderer) DocumentFooter(out *bytes.Buffer) {}

func (options *PlaintextRenderer) TocHeader(text []byte, level int) {}

func (options *PlaintextRenderer) TocFinalize() {}

func (options *PlaintextRenderer) wrap(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func (options *PlaintextRenderer) isRelativeLink(link []byte) (yes bool) {
	yes = false

	// a tag begin with '#'
	if link[0] == '#' {
		yes = true
	}

	// link begin with '/' but not '//', the second maybe a protocol relative link
	if len(link) >= 2 && link[0] == '/' && link[1] != '/' {
		yes = true
	}

	// only the root '/'
	if len(link) == 1 && link[0] == '/' {
		yes = true
	}
	return
}
