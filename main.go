package main

// Use: go build -tags "icu libstemmer"

import (
	"github.com/boombuler/knowledgedb/index"
	"github.com/boombuler/knowledgedb/log"
	"github.com/boombuler/knowledgedb/storage"
	"github.com/kardianos/service"
	"os"
	"time"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	p.exit = make(chan struct{})
	initTemplates()
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	storage.Init()

	index.Init(func() <-chan index.Document {
		result := make(chan index.Document)
		go func() {
			for e := range storage.ListAllEntries() {
				result <- e
			}
			close(result)
		}()
		return result
	})

	srv := new(httpServer)
	go srv.Serve(p.exit)
}
func (p *program) Stop(s service.Service) error {
	close(p.exit)
	time.Sleep(200 * time.Millisecond)
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "KnowledgeDB",
		DisplayName: "Knowledge DB",
		Description: "Knowledge DB - Webserver",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.SetServiceLogger(logger)

	if len(os.Args) == 2 {
		command := os.Args[1]
		err := service.Control(s, command)
		if err != nil {
			log.Fatalf("%v\n\nValidOperations: %v", err, service.ControlAction)
		}
		return
	}
	err = s.Run()
	if err != nil {
		log.Error(err)
	}
}
