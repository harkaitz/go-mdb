package mdb

import (
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"fmt"
	S "strings"
)

type Mdb interface {
	Validate()            error
	GetId()               string
	Set(key, val string)  error
	Marshal()             ([]byte, error)
	Unmarshal(b []byte)   error
}

func SetKV(o Mdb, args []string) (e error) {
	for _, arg := range args {
		key, val, splet := S.Cut(arg, "=")
		if !splet {
			return fmt.Errorf("%s: Missing value", arg)
		}
		e = o.Set(key, val)
		if (e != nil) {
			return fmt.Errorf("%s: %s", arg, e)
		}
	}
	return
}

func Save(db *redis.Client, key string, o Mdb) (err error) {
	j, err := o.Marshal()
	if err != nil {
		return
	}
	id := o.GetId()
	if len(id) == 0 {
		id = uuid.New().String()
	}
	err = o.Validate()
	if err != nil {
		return
	}
	return db.HSet(key, id, j).Err()
}

func Load(db *redis.Client, key string, id string, o Mdb) (err error) {
	j, err := db.HGet(key, id).Bytes()
	if err != nil {
		return
	}
	err = o.Unmarshal(j)
	return
}

func List(db *redis.Client, key string) (ids []string, err error) {
	return db.HKeys(key).Result()
}

func Del(db *redis.Client, key string, id string) (err error) {
	return db.HDel(key, id).Err()
}

func DelAll(db *redis.Client, key string) (err error) {
	return db.Del(key).Err()
}
