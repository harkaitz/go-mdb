package mdb

import (
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/harkaitz/go-recutils"
	"os"
	"io"
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

var DB *redis.Client

func init() {
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
	DB = redis.NewClient(&redis_opts)
	if err := DB.Ping().Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
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

func Save(key string, id string, o Mdb) (err error) {
	j, err := o.Marshal()
	if err != nil {
		return
	}
	err = o.Validate()
	if err != nil {
		return
	}
	return DB.HSet(key, id, j).Err()
}

func Load(key string, id string, o Mdb) (err error) {
	j, err := DB.HGet(key, id).Bytes()
	if err != nil {
		return
	}
	return o.Unmarshal(j)
}

func List(key string) (ids []string, err error) {
	return DB.HKeys(key).Result()
}

func Del(key string, id string) (err error) {
	return DB.HDel(key, id).Err()
}

func DelAll(key string) (err error) {
	return DB.Del(key).Err()
}

func CmdList(key string) (err error) {
	var ids []string
	ids, err = List(key)
	if err != nil {
		return err
	}
	for _, id := range ids {
		fmt.Printf("%s\n", id)
	}
	return nil
}

func CmdDel(key string, ids []string) (err error) {
	for _, id := range ids {
		err = Del(key, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func CmdAdd(key string, obj Mdb, args []string) (err error) {
	id  := uuid.New().String()

	if err = SetKV(obj, args); err != nil {
		return err
	}
		
	return Save(key, id, obj)
}

func CmdGet(key string, id string) (err error) {
	j, err := DB.HGet(key, id).Bytes()
	if err != nil {
		return
	}
	fmt.Printf("%s\n", j)
	return nil
}

func CmdSet(key string, id string, obj Mdb, args []string) (err error) {
	if err = Load(key, id, obj); err != nil {
		return err
	}
	if err = SetKV(obj, args); err != nil {
		return err
	}
	return Save(key, id, obj)
}

func CmdReadRec(key string, obj Mdb, fp *os.File) (err error) {
	var rec     *recfile.Reader
	var fields []recfile.Field
	var id     string

	rec  = recfile.NewReader(os.Stdin)

	for {
		fields, err = rec.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		obj.Reset()
		id = ""
		for _, field := range fields {
			if S.EqualFold(field.Name, "id") {
				id = field.Value
			}
			err = obj.Set(field.Name, field.Value)
			if err != nil {
				return err
			}
		}
		if len(id)==0 {
			err = fmt.Errorf("Registry without Id field.")
			return
		}
		err = Save(key, id, obj)
		if err != nil {
			return err
		}
	}
	
	return nil
}
