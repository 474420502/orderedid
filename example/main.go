package main

import (
	"log"

	"github.com/474420502/orderedid"
)

func main1() {
	var id orderedid.OrderedID
	creator := orderedid.NewWith(1)
	id = creator.Create()
	log.Println(id.Uint64(), id.Base58(), id.Timestamp(), id.NodeID()) // 142125288653825 27noD5f5R 1646623082475 1
}

func main() {
	var id orderedid.OrderedID
	creator := orderedid.New()
	id = creator.Create()
	log.Println(id.Uint64(), id.Base64(), id.Timestamp(), id.NodeID()) // 142125288653825 27noD5f5R 1646623082475 1
}
