package main

import (
	"github.com/harkaitz/go-mdb"
	"github.com/pborman/getopt/v2"
	"fmt"
	"os"
	"log"
)

const help string =
`Usage: go-mdb ...

Manage "go-mdb" objects. 

    -k REDIS_KEY   : Redis key to use.
    -l             : List object ids.
    -a     K=V ... : Add object.
    -g ID          : Get object in JSON format.
    -s ID  K=V ... : Set parameters to object.
    -d ID...       : Delete objects.
    -D             : Delete all objects in redis key.`
const copyrightLine string =
`Bug reports, feature requests to gemini|https://harkadev.com/oss
Copyright (c) 2022 Harkaitz Agirre, harkaitz.aguirre@gmail.com`

func main() {

	var err  error
	var obj  mdb.Object
	
	getopt.BoolLong("help", 'h')
	kFlag := getopt.String('k', "GO-MDB::DEFAULT")
	lFlag := getopt.Bool('l')
	aFlag := getopt.Bool('a')
	gFlag := getopt.String('g', "")
	sFlag := getopt.String('s', "")
	dFlag := getopt.Bool('d')
	DFlag := getopt.Bool('D')
	
	getopt.SetUsage(func() { fmt.Println(help + "\n\n" + copyrightLine) })
	getopt.Parse()

	key := *kFlag
	op  := ""
	id  := ""
	if *lFlag        { op += "l" }
	if *aFlag        { op += "a" }
	if len(*gFlag)>0 { op += "g"; id = *gFlag }
	if len(*sFlag)>0 { op += "s"; id = *sFlag }
	if *dFlag        { op += "d" }
	if *DFlag        { op += "D" }
	
	switch len(op) {
	case 0: getopt.Usage(); os.Exit(0);
	case 1:
	default: log.Fatal("Please specify only one operation.")
	}

	db := mdb.NewRedisClient()

	switch op {
	case "l": err = mdb.CmdList(db, key)
	case "a": err = mdb.CmdAdd(db, key, &obj, getopt.Args())
	case "g": err = mdb.CmdGet(db, key, id)
	case "s": err = mdb.CmdSet(db, key, id, &obj, getopt.Args())	
	case "d": err = mdb.CmdDel(db, key, getopt.Args())
	case "D": err = mdb.DelAll(db, key)
	}
	if err != nil {
		log.Fatal(err)
	}

}
