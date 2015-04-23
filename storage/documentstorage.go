package storage

import (
	"github.com/boltdb/bolt"
	"github.com/boombuler/knowledgedb/index"
	"github.com/boombuler/knowledgedb/log"
	"github.com/nu7hatch/gouuid"
	"time"
)

const document_dir = "documents"

var documentBucket = []byte("Documents")
var tagBucket = []byte("Tags")

func CreateEntry(entry *Entry) string {
	uid, _ := uuid.NewV4()
	id := uid.String()

	StoreEntry(id, entry)
	index.Index(id, entry.AsPlaintext())
	return id
}

func UpdateEntry(id string, entry *Entry) {
	StoreEntry(id, entry)
	index.Index(id, entry.AsPlaintext())
}

func DeleteEntry(id string, entry *Entry) {
	err := storage.Update(func(tx *bolt.Tx) error {
		for _, t := range entry.Tags {
			removeIDFromTag(tx, t, id)
		}
		b, _ := tx.CreateBucketIfNotExists(documentBucket)
		b.Delete([]byte(id))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	index.Delete(id)
}

func addIDToTag(tx *bolt.Tx, tag, id string) error {
	buck, err := tx.CreateBucketIfNotExists(tagBucket)
	if err != nil {
		return err
	}
	buck, err = buck.CreateBucketIfNotExists([]byte(tag))
	if err != nil {
		return err
	}
	buck.Put([]byte(id), []byte{1})
	return nil
}

func removeIDFromTag(tx *bolt.Tx, tag, id string) error {
	buck, err := tx.CreateBucketIfNotExists(tagBucket)
	if err != nil {
		return err
	}
	b := buck.Bucket([]byte(tag))
	if b == nil {
		return nil
	}
	if err := b.Delete([]byte(id)); err != nil {
		if err == bolt.ErrIncompatibleValue {
			return nil
		}
		return err
	}

	if k, _ := b.Cursor().Next(); k == nil {
		return buck.DeleteBucket([]byte(tag))
	}
	return nil
}

func containsTag(tag string, tagList []string) bool {
	if tagList == nil {
		return false
	}
	for _, t := range tagList {
		if t == tag {
			return true
		}
	}
	return false
}

func StoreEntry(id string, entry *Entry) {
	err := storage.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(documentBucket)
		if err != nil {
			return err
		}
		oldTags, err := fetchOldTags(b, id)
		if err != nil {
			return err
		}

		for _, t := range oldTags {
			if !containsTag(t, entry.Tags) {
				if err := removeIDFromTag(tx, t, id); err != nil {
					return err
				}
			}
		}
		for _, t := range entry.Tags {
			if !containsTag(t, oldTags) {
				if err := addIDToTag(tx, t, id); err != nil {
					return err
				}
			}
		}

		entry.LastModified = time.Now()
		return b.Put([]byte(id), encodeData(entry))
	})
	if err != nil {
		log.Fatal(err)
	}
}

func fetchOldTags(docBucket *bolt.Bucket, id string) ([]string, error) {
	entry := new(Entry)
	decodeData(docBucket.Get([]byte(id)), entry)
	return entry.Tags, nil
}

func LoadEntry(id string) (entry *Entry) {
	err := storage.View(func(tx *bolt.Tx) error {
		entry = new(Entry)
		b := tx.Bucket(documentBucket)
		if b == nil {
			return nil
		}
		decodeData(b.Get([]byte(id)), entry)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return entry
}

func ListAllEntries() <-chan *EntryWithId {
	result := make(chan *EntryWithId)
	go storage.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(documentBucket)
		if bucket != nil {
			c := bucket.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if v != nil {
					id := string(k)
					entry := new(Entry)
					decodeData(v, entry)
					result <- &EntryWithId{id, *entry}
				}
			}
		}
		close(result)
		return nil
	})
	return result
}
