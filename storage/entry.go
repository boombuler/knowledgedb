package storage

import (
	"github.com/boombuler/knowledgedb/markdown"
	"time"
)

type Entry struct {
	Title        string
	Body         string
	Tags         []string
	LastModified time.Time
}

func (e *Entry) Type() string {
	return "entry"
}

func (e *Entry) UnifyTags() {
	var tagsNew = make([]string, 0)
	for _, t := range e.Tags {
		if !containsTag(t, tagsNew) {
			tagsNew = append(tagsNew, t)
		}
	}
	e.Tags = tagsNew
}

func (e Entry) AsPlaintext() *Entry {
	bodyNew := markdown.ToPlaintext(e.Body)
	return &Entry{
		Title: e.Title,
		Body:  bodyNew,
		Tags:  e.Tags,
	}
}

const body_snippet_size = 200

func (e Entry) GetBodySnippet() string {
	body := e.AsPlaintext().Body
	length := StrLen(body)
	if length > body_snippet_size {
		return SubStr(body, 0, body_snippet_size)
	}
	return body
}

type EntryWithId struct {
	id string
	Entry
}

func NewEntryWithId(id string, entry *Entry) *EntryWithId {
	return &EntryWithId{
		id,
		*entry,
	}
}

func (e *EntryWithId) Id() string {
	return e.id
}

func (e *EntryWithId) Content() interface{} {
	return e.AsPlaintext()
}
