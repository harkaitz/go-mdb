package mdb

import (
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"os"
	"fmt"
	S "strings"
)

type Mdb interface {
	Reset()
	Validate()            error
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

func Save(db *redis.Client, key string, id string, o Mdb) (err error) {
	j, err := o.Marshal()
	if err != nil {
		return
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
	o.Reset()
	return o.Unmarshal(j)
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

func NewRedisClient() (db *redis.Client) {
	redis_address  := os.Getenv("REDIS_ADDRESS")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_opts     := redis.Options {}
	if len(redis_address)>0 {
		redis_opts.Addr = redis_address
	} else {
		redis_opts.Addr = "localhost:6379"
	}
	if len(redis_password)>0 {
		redis_opts.Password = redis_password
	}
	return redis.NewClient(&redis_opts)
}

func CmdList(db *redis.Client, key string) (err error) {
	var ids []string
	ids, err = List(db, key)
	if err != nil {
		return err
	}
	for _, id := range ids {
		fmt.Printf("%s\n", id)
	}
	return nil
}

func CmdDel(db *redis.Client, key string, ids []string) (err error) {
	for _, id := range ids {
		err = Del(db, key, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func CmdAdd(db *redis.Client, key string, obj Mdb, args []string) (err error) {
	id  := uuid.New().String()

	if err = SetKV(obj, args); err != nil {
		return err
	}
		
	return Save(db, key, id, obj)
}

func CmdGet(db *redis.Client, key string, id string) (err error) {
	j, err := db.HGet(key, id).Bytes()
	if err != nil {
		return
	}
	fmt.Printf("%s\n", j)
	return nil
}

func CmdSet(db *redis.Client, key string, id string, obj Mdb, args []string) (err error) {
	if err = Load(db, key, id, obj); err != nil {
		return err
	}
	if err = SetKV(obj, args); err != nil {
		return err
	}
	return Save(db, key, id, obj)
}
