package mqEngin

import (
	"encoding/binary"
	"fmt"
	"github.com/jacoblai/httprouter"
	"github.com/jacoblai/yiyidb"
	"io/ioutil"
	"net/http"
	"snowflake"
	"strconv"
	"sync"
)

type MqEngin struct {
	localdb  *yiyidb.Kvdb
	id       uint64
	idWroker *snowflake.Worker
	pool     *sync.Pool
}

func NewMqEngin(dir string) *MqEngin {
	worker, err := snowflake.NewWorker(1)
	if err != nil {
		panic(err)
	}
	km := &MqEngin{
		idWroker: worker,
		pool: &sync.Pool{
			New: func() interface{} {
				cc := make([]byte, 8)
				return cc
			},
		},
	}
	//初始化id
	km.id = 0
	db, err := yiyidb.OpenKvdb(dir+"/db", false, true, 8)
	if err != nil {
		panic(err)
	}
	//所有数据库key以int64实现，因此数据库会自动排序，从小到大。
	//取库已存在数据的最大值作为重启后的最大id，防止id错乱
	k, err := db.GetLastKey()
	if err == nil {
		id := km.KeyToId(k)
		if id > km.id {
			km.id = id
		}
	}
	km.localdb = db
	return km
}

func (k *MqEngin) Close() {
	k.localdb.Close()
}

func (k *MqEngin) IdToKey(id uint64) []byte {
	key := k.pool.Get().([]byte)
	defer k.pool.Put(key)
	binary.BigEndian.PutUint64(key, id)
	return key
}

func (k *MqEngin) KeyToId(key []byte) uint64 {
	return binary.BigEndian.Uint64(key)
}

func (k *MqEngin) DelQeueu(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ids := ps.ByName("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = k.localdb.Del(yiyidb.IdToKeyPure(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = fmt.Fprintf(w, `{"ok":%v}`, true)
}

func (k *MqEngin) PeekQeueu(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ids := ps.ByName("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val, err := k.localdb.Get(yiyidb.IdToKeyPure(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = fmt.Fprintf(w, `{"data":%s}`, val)
}

func (k *MqEngin) EnQeueu(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	tl := r.URL.Query().Get("ttl")
	ttl, err := strconv.Atoi(tl)
	if err != nil {
		ttl = 0
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nid := k.idWroker.GetId()
	err = k.localdb.Put(yiyidb.IdToKeyPure(nid), body, ttl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, `{"id":%v}`, nid)
}
