package main

import "log"

func main() {
	store, err := newPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	err1 := store.InitDb()
	if err != nil {
		log.Fatal(err1)
	}
	apiServer := NewApiServer(":3000", store)
	apiServer.run()
}
