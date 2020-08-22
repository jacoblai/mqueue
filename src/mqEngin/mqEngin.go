package mqEngin

import (
	"fmt"
	"github.com/jacoblai/httprouter"
	"io/ioutil"
	"net/http"
	"snowflake"
	"utils"
)

type MqEngin struct {
	localdb  *Queue
	idWroker *snowflake.Worker
	ut       *utils.Utils
}

func NewMqEngin(dir string) *MqEngin {
	worker, err := snowflake.NewWorker(1)
	if err != nil {
		panic(err)
	}
	km := &MqEngin{
		idWroker: worker,
		ut:       utils.NewUtils(),
	}
	db, err := OpenQueue(dir + "/mq")
	if err != nil {
		panic(err)
	}
	km.localdb = db
	return km
}

func (k *MqEngin) Close() {
	k.localdb.Close()
}

func (k *MqEngin) PeekQeueu(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	key, val, err := k.localdb.Peek()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, `{"id":%v,"data":%s}`, key, val)
}

func (k *MqEngin) DeQeueu(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	key, val, err := k.localdb.Dequeue()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, _ = fmt.Fprintf(w, `{"id":%v,"data":%s}`, key, val)
}

func (k *MqEngin) EnQeueu(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := k.localdb.Enqueue(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, `{"id":%v}`, id)
}
