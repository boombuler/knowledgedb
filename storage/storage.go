package storage

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	c "github.com/boombuler/knowledgedb/config"
	"path"
)

var storage *bolt.DB

func Init() {
	var err error
	storage, err = bolt.Open(path.Join(c.Current.DataDir, "documents.bolt"), 0600, nil)
	if err != nil {
		panic(err)
	}
	gob.Register(new(Entry))
}

func encodeData(obj interface{}) []byte {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	enc.Encode(obj)
	return buff.Bytes()
}

func decodeData(data []byte, out interface{}) {
	if data != nil && len(data) > 0 {
		buff := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buff)
		dec.Decode(out)
	}
}

type EntriesByDateTime []*EntryWithId

func (a EntriesByDateTime) Len() int           { return len(a) }
func (a EntriesByDateTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntriesByDateTime) Less(i, j int) bool { return a[i].LastModified.After(a[j].LastModified) }
