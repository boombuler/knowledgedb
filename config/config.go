package config

import (
	"encoding/json"
	"github.com/boombuler/knowledgedb/log"
	"os"
	"path"
)

type Config struct {
	HttpAddr        string
	DataDir         string
	DefaultAnalyzer string
}

var (
	Current        Config
	ApplicationDir string
)

func init() {
	ApplicationDir = path.Dir(os.Args[0])

	Current = Config{
		HttpAddr:        ":80",
		DataDir:         ApplicationDir,
		DefaultAnalyzer: "standard",
	}

	f, err := os.Open(path.Join(ApplicationDir, "config.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Current)
	if err != nil {
		log.Fatal(err)
	}
}
