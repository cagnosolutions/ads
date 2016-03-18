package main

import (
	"log"

	"github.com/cagnosolutions/ads/tmp"
)

type User struct {
	Id     int      `json:"id,omitempty"`
	Name   []string `json:"name,omitempty"`
	Active bool     `json:"active,omitempty"`
}

func main() {
	ndb := tmp.NewDB("db")
	err := ndb.AddStore("user")
	if err != nil {
		log.Println(err)
	}
	err = ndb.Add("user", []byte(`doc-0`), User{0, []string{"scott", "cagno"}, true})
	if err != nil {
		log.Println(err)
	}
	err = ndb.Add("user", []byte(`doc-1`), User{1, []string{"kayla", "cagno"}, false})
	if err != nil {
		log.Println(err)
	}
	err = ndb.Add("user", []byte(`doc-2`), User{2, []string{"gabe", "witmer"}, true})
	if err != nil {
		log.Println(err)
	}
	err = ndb.Add("user", []byte(`doc-3`), User{3, []string{"greg", "pechiro"}, true})
	if err != nil {
		log.Println(err)
	}
	err = ndb.Add("user", []byte(`doc-4`), User{4, []string{"rosalie", "pechiro"}, false})
	if err != nil {
		log.Println(err)
	}
}
