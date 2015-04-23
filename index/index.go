package index

import (
	"github.com/blevesearch/bleve"
	"github.com/boombuler/knowledgedb/config"
	"github.com/boombuler/knowledgedb/log"

	"path"
)

const index_name = "knowdb.bleve"
const tag_analyzer = "keyword"

var index bleve.Index

type Document interface {
	Id() string
	Content() interface{}
}

type GetAllDocsFn func() <-chan Document

func Init(allDocs GetAllDocsFn) {
	var err error

	idx_path := path.Join(config.Current.DataDir, index_name)

	index, err = bleve.Open(idx_path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Info("Creating new Index")
		indexMapping := bleve.NewIndexMapping()
		indexMapping.DefaultAnalyzer = config.Current.DefaultAnalyzer

		entryMapping := bleve.NewDocumentMapping()

		textField := bleve.NewTextFieldMapping()

		entryMapping.AddFieldMappingsAt("Body", textField)
		entryMapping.AddFieldMappingsAt("Title", textField)

		tagField := bleve.NewTextFieldMapping()
		tagField.Analyzer = tag_analyzer

		entryMapping.AddFieldMappingsAt("Tags", tagField)

		indexMapping.AddDocumentMapping("entry", entryMapping)

		index, err = bleve.New(idx_path, indexMapping)
		if err != nil {
			log.Fatal(err)
		}

		// reindex existing documents
		indexRebuildingLogged := false
		for itm := range allDocs() {
			if !indexRebuildingLogged {
				indexRebuildingLogged = true
				log.Info("Start rebuilding Search-Index")
			}

			index.Index(itm.Id(), itm.Content())
		}
		if indexRebuildingLogged {
			log.Info("Finished rebuilding Search-Index")
		}

	} else if err == nil {
		log.Info("Opening existing Index")
	} else {
		log.Fatal(err)
	}
}

func Index(id string, content interface{}) {
	index.Index(id, content)
}

func Delete(id string) {
	index.Delete(id)
}

func Search(q string) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(q)
	req := bleve.NewSearchRequest(query)
	req.Highlight = bleve.NewHighlightWithStyle("html_ex")
	return index.Search(req)
}
