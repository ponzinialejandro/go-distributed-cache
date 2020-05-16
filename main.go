package main

import (
	cache "git.topfreegames.com/alejandro.ponzini1/go-distributed-cache/cache"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	rwMutex := sync.RWMutex{}
	lruCache := cache.NewLRUCache(10000)
	http.HandleFunc("/", Hanlder{rwMutex: &rwMutex, cache: lruCache}.handle)
	http.ListenAndServe(":8080", nil)
}

type Hanlder struct {
	rwMutex *sync.RWMutex
	cache   cache.LRUCache
}

func (p Hanlder) handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPut:
		p.PutKey(w, r)
	case http.MethodGet:
		p.GetValue(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (p Hanlder) PutKey(w http.ResponseWriter, r *http.Request) {

	if id, ok := r.URL.Query()["id"]; ok {
		if bytes, err := ioutil.ReadAll(r.Body); err == nil {
			p.rwMutex.Lock()
			defer p.rwMutex.Unlock()
			p.cache.Put(id[0], string(bytes))
		}
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (p Hanlder) GetValue(w http.ResponseWriter, r *http.Request) {
	if id, ok := r.URL.Query()["id"]; ok {
		p.rwMutex.RLock()
		defer p.rwMutex.RUnlock()
		if value, found := p.cache.Get(id[0]); found {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
