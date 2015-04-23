# KnowledgeDB [![License: Apache](https://img.shields.io/:license-Apache_2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

This is a little Knowledge Database for personal use written in Go
There is no user authentication just some markdown processing and full-text-search.

## Used 3rd Party Packages / Libs

* [Blackfriday](https://github.com/russross/blackfriday) for markdown processing
* [bluemonday](https://github.com/microcosm-cc/bluemonday) for HTML sanitizing
* [bleve](https://github.com/blevesearch/bleve) for searching within the documents
* [Bolt](https://github.com/boltdb/bolt) for document storage
* [Bootstrap](http://getbootstrap.com/) basic UI design
* [gorilla/mux](https://github.com/gorilla/mux) for http request routing
* [highlight.js](https://highlightjs.org/) for syntax highlighting
* [JQuery](https://jquery.com/) because some other libs needed it...
* [Bootstrap Tags Input](https://github.com/timschlechter/bootstrap-tagsinput/) for tag input :-P
* [typeahead.js](https://twitter.github.io/typeahead.js/) also for tag input
* [gouuid](https://github.com/nu7hatch/gouuid) to generate unique IDs.
* [service](github.com/kardianos/service) for running the webserver as a windows service.

## Compile
To compile you could simply use `go build` or if you want some extra analyzers for the search engine build with `go build -tags "icu libstemmer"`
For more information about additional search engine features see [bleve](https://github.com/blevesearch/bleve)...

## Executing
Copy the compiled executable in the same directory as the `static` and the `templates` directory.
If you want you can create a `config.json`

```json
{
    "HttpAddr": ":12345",               // port / address for the http server
    "DataDir": "/var/knowledgedb/data", // directory to store the database and fulltext-search-index
    "DefaultAnalyzer": "standard"       // Textanalyzer for the bleve search.
}
```

