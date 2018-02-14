package adapters

import (
	"github.com/mediocregopher/radix.v2/redis"
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

	slog.Infof("[storage.redis -> Connect] Connecting: %s://%s\n", f.network, f.addr)

	if !f.isConnected(){

		if conn, err := redis.Dial(f.network, f.addr) ; err == nil{
			f.storage = conn
		}else{
			slog.Panicln(err)
		}
	}

	isConnected = f.isConnected()

	var version string
	if version = f.Version(); version == "0" {version = f.incVersion()}

	slog.Infoln("[storage.redis -> Connect] Connected: ", f.isConnected())
	slog.Infoln("[storage.redis -> Connect] Version: ", version)

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

	slog.Debugf("[storage.redis -> Add] In: \n\t-- index: `%s` \n\t-- record: `%s` \n\t-- force: %v\n", index, record, force)

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

	slog.Debugln("[storage.redis -> Add] ", "error: ", err)

	return err
}

func (f *storageRedis) Rm(index string) bool{

	defer f.UnLock(index)

	var r int
	if f.Lock(index){

		r, _ = f.storage.Cmd("HDEL", f.storageKey, index).Int()

		f.incVersion()

	}else{
		slog.Infoln("[storage.redis -> Add] Pull: lock fail for", index)
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

	slog.Debugln("[storage.redis -> incVersion] Version: ", "update:", oldVersion,"-->",intVersion)

	return version
}


func (f *storageRedis) Lock(index string) bool{

	if index == ""{
		slog.Panicln("index not specified")
	}

	slog.Debugf("[storage.redis -> Lock] Command: HSETNX %s %s %d\n", f.storageLock, index, 1)

	l,e := f.storage.Cmd("HSETNX", f.storageLock, index, 1).Int()

	slog.Debugln("[storage.redis -> Lock] Result: ", l,e)

	return l==1
}

func (f *storageRedis) UnLock(index string) bool {

	slog.Debugln("[storage.redis -> UnLock] Command: ", "HDEL", f.storageLock, index)

	r,e := f.storage.Cmd("HDEL", f.storageLock, index).Int()
	if e != nil{panic(e)}

	return r == 1
}

func (f *storageRedis) pull(index string) (record string){

	defer f.UnLock(index)

	slog.Infoln("[storage.redis] Pull: ", index)
	if f.Lock(index){

		if record,_ = f.storage.Cmd("HGET", f.storageKey, index).Str(); record == ""{
			slog.Infoln("[storage.redis -> pull] Pull: no data for index", index)
		}

		f.storage.Cmd("HDEL", f.storageKey, index)

		f.incVersion()

	}else{
		slog.Infoln("[storage.redis -> pull] Pull: lock fail for", index)
	}

	return record
}
