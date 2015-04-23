package storage

import (
	"github.com/boltdb/bolt"
)

type TagInfo struct {
	Name  string
	Count int
}

type TagInfos []TagInfo

func (a TagInfos) Len() int           { return len(a) }
func (a TagInfos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TagInfos) Less(i, j int) bool { return StringIsLessIgnoreCase(a[i].Name, a[j].Name) }

func GetAllTags() TagInfos {
	result := make(TagInfos, 0)
	storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tagBucket)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil { // Sub bucket
				tag := b.Bucket(k)
				if tag != nil {
					val := string(k)
					result = append(result, TagInfo{
						val,
						tag.Stats().KeyN,
					})
				}
			}
		}
		return nil
	})

	return result
}

func GetAllEntriesWithTag(tag string) []*EntryWithId {
	result := make([]*EntryWithId, 0)
	storage.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tagBucket)
		if b == nil {
			return nil
		}
		b = b.Bucket([]byte(tag))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				entry := LoadEntry(string(k))
				if entry != nil {
					result = append(result, NewEntryWithId(string(k), entry))
				}
			}
		}
		return nil
	})

	return result
}
