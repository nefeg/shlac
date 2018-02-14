package adapters

import (
	"github.com/mediocregopher/radix.v2/redis"
	"log"
	"strconv"
	"errors"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

var errIndexExist = errors.New("index already exist")

type storageRedis struct {

	network     string
	addr        string
	storageKey  string
	storageLock string
	storageVer  string

	storage     *redis.Client
}

func NewRedisAdapter(network, addr, storageKey string) *storageRedis{

	s := &storageRedis{network:network, addr:addr}
	s.storageKey    = storageKey + ".db"
	s.storageLock   = s.storageKey + ".lock"
	s.storageVer    = s.storageKey + ".version"

	s.Connect()

	return s
}

func (f *storageRedis) Connect() (isConnected bool){

	slog.InfoF("[storage.redis -> Connect] Connecting: %s://%s\n", f.network, f.addr)

	if !f.isConnected(){

		if conn, err := redis.Dial(f.network, f.addr) ; err == nil{
			f.storage = conn
		}else{
			log.Panicln(err)
		}
	}

	isConnected = f.isConnected()

	var version string
	if version = f.Version(); version == "0" {version = f.incVersion()}

	slog.InfoLn("[storage.redis -> Connect] Connected: ", f.isConnected())
	slog.InfoLn("[storage.redis -> Connect] Version: ", version)

	return isConnected
}

func (f *storageRedis) Disconnect(){
	f.storage.Close()
}

func (f *storageRedis) isConnected() bool{
	return f.storage != nil
}



func (f *storageRedis) Exists(index string) bool{

	resp, _ := f.storage.Cmd("HEXISTS", f.storageKey, index).Int()

	return resp>0
}

func (f *storageRedis) Get(index string) (record string){

	record,_ = f.storage.Cmd("HGET", f.storageKey, index).Str()

	return record
}

func (f *storageRedis) Add(index string, record string, force bool) (err error){

	slog.DebugF("[storage.redis -> Add] In: \n\t-- index: `%s` \n\t-- record: `%s` \n\t-- force: %v\n", index, record, force)

	var resp int

	if force {
		resp, err = f.storage.Cmd("HSET", f.storageKey, index, record).Int()
	}else{
		resp, err = f.storage.Cmd("HSETNX", f.storageKey, index, record).Int()
	}

	if !(resp > 0) {
		err = errIndexExist

	}else{
		f.incVersion()
	}

	slog.DebugLn("[storage.redis -> Add] ", "error: ", err)

	return err
}

func (f *storageRedis) Rm(index string) bool{

	defer f.UnLock(index)

	var r int
	if f.Lock(index){

		r, _ = f.storage.Cmd("HDEL", f.storageKey, index).Int()

		f.incVersion()

	}else{
		log.Println("[storage.redis -> Add] Pull: lock fail for", index)
	}


	return r>0
}

func (f *storageRedis) List() (data map[string]string){

	data, _ = f.storage.Cmd("HGETALL", f.storageKey).Map()

	return data
}

func (f *storageRedis) Flush(){
	f.incVersion()
	f.storage.Cmd("DEL", f.storageKey)
}


func (f *storageRedis) Version() (version string){

	intVersion, _ := f.storage.Cmd("GET", f.storageVer).Int()

	version = strconv.Itoa(intVersion)

	return version
}

func (f *storageRedis) incVersion() (version string){

	oldVersion := f.Version()

	intVersion, _ := f.storage.Cmd("INCR", f.storageVer).Int()

	version = strconv.Itoa(intVersion)

	slog.DebugLn("[storage.redis  -> incVersion] Version: ", "update:", oldVersion,"-->",intVersion)

	return version
}


func (f *storageRedis) Lock(index string) bool{

	if index == ""{
		slog.PanicLn("index not specified")
	}

	slog.DebugF("[storage.redis -> Lock] Command: HSETNX %s %s %d\n", f.storageLock, index, 1)

	l,e := f.storage.Cmd("HSETNX", f.storageLock, index, 1).Int()

	slog.DebugLn("[storage.redis -> Lock] Result: ", l==1,e)

	return l==1
}

func (f *storageRedis) UnLock(index string) {

	slog.DebugLn("[storage.redis  -> UnLock] Command: ", "HDEL", f.storageLock, index)

	f.storage.Cmd("HDEL", f.storageLock, index)
}

func (f *storageRedis) pull(index string) (record string){

	defer f.UnLock(index)

	log.Println("[storage.redis] Pull: ", index)
	if f.Lock(index){

		if record,_ = f.storage.Cmd("HGET", f.storageKey, index).Str(); record == ""{
			log.Println("[storage.redis -> pull] Pull: no data for index", index)
		}

		f.storage.Cmd("HDEL", f.storageKey, index)

		f.incVersion()

	}else{
		log.Println("[storage.redis -> pull] Pull: lock fail for", index)
	}

	return record
}
