package main

import (
	"github.com/boombuler/knowledgedb/config"
	"github.com/boombuler/knowledgedb/log"
	"github.com/boombuler/knowledgedb/storage"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"path"
	"sort"
	"time"
)

type httpServer struct{}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (s httpServer) Serve(exit <-chan struct{}) {
	r := mux.NewRouter()
	r.HandleFunc("/new", s.serveCreateEntry).Methods("GET", "POST")
	r.HandleFunc("/delete/{id}", s.withDoc(s.serveDeleteEntry)).Methods("POST")
	r.HandleFunc("/edit/{id}", s.withDoc(s.serveEditEntry)).Methods("GET", "POST")
	r.HandleFunc("/view/{id}", s.withDoc(s.serveViewEntry)).Methods("GET")

	r.HandleFunc("/search", s.serveSearchEntries).Methods("POST")

	r.HandleFunc("/alltags", s.serveAllTags).Methods("GET")
	r.HandleFunc("/tag/{tag}", s.serveTagDetails).Methods("GET")
	r.HandleFunc("/licenses", s.serveLicenses).Methods("GET")
	r.HandleFunc("/help", s.serveMarkdown("help")).Methods("GET")
	r.HandleFunc("/", s.serveIndex).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(path.Join(config.ApplicationDir, "static")))).Methods("GET")

	ln, err := net.Listen("tcp", config.Current.HttpAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	srv := http.Server{Handler: r}
	go func() {
		log.Info("Starting Webserver")
		srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	}()
	select {
	case <-exit:
	}
	log.Info("Stopping Webserver")
	ln.Close()
}

func (s *httpServer) serveIndex(w http.ResponseWriter, r *http.Request) {
	tags := storage.GetAllTags()
	sort.Sort(tags)
	renderTemplate(w, "index", struct {
		Tags storage.TagInfos
	}{
		tags,
	})
}

func (s *httpServer) serveMarkdown(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		renderMarkdown(w, "help")
	}
}
