package mqEngin

import (
	"errors"
	"github.com/jacoblai/yiyidb"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"snowflake"
)

const (
	KB int = 1024
	MB int = KB * 1024
	GB int = MB * 1024
)

var (
	ErrEmpty       = errors.New("queue is empty")
	ErrOutOfBounds = errors.New("ID used is outside range of queue")
	ErrDBClosed    = errors.New("Database is closed")
)

//FIFO
type Queue struct {
	DataDir      string
	db           *leveldb.DB
	iteratorOpts *opt.ReadOptions
	idWorker     *snowflake.Worker
}

func OpenQueue(dataDir string) (*Queue, error) {
	var err error
	wk, err := snowflake.NewWorker(1)
	q := &Queue{
		DataDir:      dataDir,
		db:           &leveldb.DB{},
		iteratorOpts: &opt.ReadOptions{DontFillCache: true},
		idWorker:     wk,
	}

	opts := &opt.Options{}
	opts.ErrorIfMissing = false
	opts.BlockCacheCapacity = 4 * MB
	//队列key固定用8个byte所以bloom应该是8*1.44~12优化查询
	opts.Filter = filter.NewBloomFilter(12)
	opts.Compression = opt.SnappyCompression
	opts.BlockSize = 4 * KB
	opts.WriteBuffer = 4 * MB
	opts.OpenFilesCacheCapacity = 1 * KB
	opts.CompactionTableSize = 32 * MB
	opts.WriteL0SlowdownTrigger = 16
	opts.WriteL0PauseTrigger = 64

	q.db, err = leveldb.OpenFile(dataDir, opts)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) EnqueueBatch(value [][]byte) error {
	batch := new(leveldb.Batch)
	for _, v := range value {
		id := q.idWorker.GetId()
		batch.Put(yiyidb.IdToKeyPure(id), v)
	}
	if err := q.db.Write(batch, nil); err != nil {
		return err
	}

	return nil
}

func (q *Queue) Enqueue(value []byte) (int64, error) {
	id := q.idWorker.GetId()
	if err := q.db.Put(yiyidb.IdToKeyPure(id), value, nil); err != nil {
		return 0, err
	}
	return id, nil
}

func (q *Queue) Dequeue() (int64, []byte, error) {
	iter := q.db.NewIterator(nil, q.iteratorOpts)
	defer iter.Release()
	if ok := iter.First(); !ok {
		return 0, nil, ErrEmpty
	}
	_ = q.db.Delete(iter.Key(), nil)

	return yiyidb.KeyToIDPure(iter.Key()), iter.Value(), nil
}

func (q *Queue) Peek() (int64, []byte, error) {
	iter := q.db.NewIterator(nil, q.iteratorOpts)
	defer iter.Release()
	if ok := iter.First(); !ok {
		return 0, nil, ErrEmpty
	}
	return yiyidb.KeyToIDPure(iter.Key()), iter.Value(), nil
}

func (q *Queue) Update(id int64, newValue []byte) error {
	if err := q.db.Put(yiyidb.IdToKeyPure(id), newValue, nil); err != nil {
		return err
	}

	return nil
}

func (q *Queue) Close() {
	q.db.Close()
}
