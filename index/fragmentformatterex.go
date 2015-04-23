package index

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/registry"
	"github.com/blevesearch/bleve/search/highlight"
	tt "text/template"
)

type FragmentFormatterEx struct {
	before string
	after  string
}

func NewFragmentFormatterEx(before, after string) *FragmentFormatterEx {
	return &FragmentFormatterEx{
		before: before,
		after:  after,
	}
}

func (a *FragmentFormatterEx) Format(f *highlight.Fragment, orderedTermLocations highlight.TermLocations) string {
	rv := ""
	curr := f.Start
	for _, termLocation := range orderedTermLocations {
		if termLocation == nil {
			continue
		}
		if termLocation.Start < curr {
			continue
		}
		if termLocation.End > f.End {
			break
		}
		// add the stuff before this location
		rv += tt.HTMLEscapeString(string(f.Orig[curr:termLocation.Start]))
		// add the color
		rv += a.before
		// add the term itself
		rv += tt.HTMLEscapeString(string(f.Orig[termLocation.Start:termLocation.End]))
		// reset the color
		rv += a.after
		// update current
		curr = termLocation.End
	}
	// add any remaining text after the last token
	rv += tt.HTMLEscapeString(string(f.Orig[curr:f.End]))

	return rv
}

func init() {
	const FormatterName = "html_ex"
	registry.RegisterFragmentFormatter(FormatterName, func(config map[string]interface{}, cache *registry.Cache) (highlight.FragmentFormatter, error) {
		return NewFragmentFormatterEx("<span class=\"highlight\">", "</span>"), nil
	})

	_, err := bleve.Config.Cache.DefineFragmentFormatter(FormatterName,
		map[string]interface{}{
			"type": FormatterName,
		})
	if err != nil {
		panic(err)
	}

	_, err = bleve.Config.Cache.DefineHighlighter(FormatterName,
		map[string]interface{}{
			"type":       "simple",
			"fragmenter": "simple",
			"formatter":  FormatterName,
		})
	if err != nil {
		panic(err)
	}
}
