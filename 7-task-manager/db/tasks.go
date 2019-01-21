package db

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"time"
)

const (
	defaultMode             = 0600
	taskBucketName          = "tasks"
	completedTaskBucketName = "tasks_completed"
)

var store TaskDb

func Init(dbPath string) error {
	db, err := bolt.Open(dbPath, defaultMode, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return err
	}

	store = TaskDb{db}

	if err := store.CreateBucket(TaskBucket(taskBucketName)); err != nil {
		return err
	}
	if err := store.CreateBucket(TaskBucket(completedTaskBucketName)); err != nil {
		return err
	}

	return nil
}

func (store TaskDb) CreateBucket(b Bucket) error {
	if err := store.bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b.Name())
		return err
	}); err != nil {
		return err
	}
	return nil
}

type TaskManager interface {
	Add(task Task) (int, error)
	GetAll() ([]Task, error)
	Delete(key int) error
}

type TaskDb struct {
	bolt *bolt.DB
}

type Task struct {
	Key   int
	Value string
}

type CompletedTask struct {
	task     Task
	datetime int64
}

type Bucket interface {
	Name() []byte
}

type TaskBucket string

func (bucket TaskBucket) Name() []byte {
	return []byte(bucket)
}

func Add(task Task) (int, error) {
	var id int
	err := store.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName())
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task.Value))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetCompleted() ([]CompletedTask, error) {
	var tasks []CompletedTask
	err := store.bolt.View(func(tx *bolt.Tx) error {
		cB := tx.Bucket(completedBucketName())

		c := cB.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, CompletedTask{
				task: Task{
					Key:   btoi(k),
					Value: string(v),
				},
				datetime: time.Unix(int64(btoi(k)), 0).Unix(),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetAll() ([]Task, error) {
	var tasks []Task
	err := store.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName())

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func Delete(key int) error {
	return store.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName())
		return b.Delete(itob(key))
	})
}

func Complete(task Task) error {
	return store.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName())

		if err := b.Delete(itob(task.Key)); err != nil {
			return err
		}

		cB := tx.Bucket(completedBucketName())

		completedT := CompletedTask{task: task, datetime: time.Now().Unix()}

		key := itob(int(completedT.datetime))

		return cB.Put(key, []byte(completedT.task.Value))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func bucketName() []byte {
	return (TaskBucket(taskBucketName)).Name()
}

func completedBucketName() []byte {
	return (TaskBucket(taskBucketName + "_completed")).Name()
}
