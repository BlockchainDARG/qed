// Copyright © 2018 Banco Bilbao Vizcaya Argentaria S.A.  All rights reserved.
// Use of this source code is governed by an Apache 2 License
// that can be found in the LICENSE file

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"verifiabledata/api/apihttp"
	"verifiabledata/balloon"
	"verifiabledata/balloon/hashing"
	"verifiabledata/balloon/storage"
	"verifiabledata/balloon/storage/badger"
	"verifiabledata/balloon/storage/bolt"
	"verifiabledata/balloon/storage/cache"
)

var (
	httpEndpoint, dbPath, storageName string
	cacheSize                         int
)

func main() {
	// We use the TypeVar flag syntax becouse balloon requires parameters as *type
	flag.StringVar(&httpEndpoint, "http_endpoint", ":8080", "Endpoint for REST requests on (host:port)")
	flag.StringVar(&dbPath, "path", "/tmp/balloon.db", "Set default storage path.")
	flag.IntVar(&cacheSize, "cache", 5000000, "Initialize and reserve custom cache size.")
	flag.StringVar(&storageName, "storage", "badger", "Choose between different storage backends. Eg badge|bolt")
	flag.Parse()

	var frozen, leaves storage.Store

	switch storageName {
	case "badger":
		frozen = badger.NewBadgerStorage(fmt.Sprintf("%s/frozen.db", dbPath))
		leaves = badger.NewBadgerStorage(fmt.Sprintf("%s/leaves.db", dbPath))
	case "bolt":
		frozen = bolt.NewBoltStorage(fmt.Sprintf("%s/frozen.db", dbPath), "frozen")
		leaves = bolt.NewBoltStorage(fmt.Sprintf("%s/leaves.db", dbPath), "leaves")
	default:
		fmt.Print("Please select a valid storage backend")
	}

	cache := cache.NewSimpleCache(cacheSize)

	balloon := balloon.NewHyperBalloon(dbPath, hashing.Sha256Hasher, frozen, leaves, cache)

	err := http.ListenAndServe(httpEndpoint, apihttp.NewApiHttp(balloon))
	if err != nil {
		log.Fatalln("Can't start HTTP Server: ", err)
	}
}