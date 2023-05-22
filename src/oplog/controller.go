package oplog

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Controller struct {
		// Controller is the sole entity responsible for coordinating
		// between the watcher, buffer and managing the bookmarks.
		// It runs in the scope of the collection.
		Database   *mongo.Database
		Collection *mongo.Collection

		watcher *OplogWatcher
		buffer  *Buffer
	}
)

func NewController(db *mongo.Database, collection *mongo.Collection) (ctrlr *Controller, err error) {
	ctrlr = &Controller{
		Database:   db,
		Collection: collection,
	}
	if ctrlr.watcher, err = NewOplogWatcher(db, collection); err != nil {
		return
	}
	if ctrlr.buffer, err = NewBuffer(LogFlusherFunc); err != nil {
		return
	}
	return
}

func (ctrlr *Controller) trackWatcherMessages() (err error) {
	for {
		select {
		case msg := <-ctrlr.watcher.CtrlrCh:
			log.Println("Message received in controller", msg)
			if err = ctrlr.buffer.Store(msg); err != nil {
				log.Println("Error on storing message in buffer", msg, err)
			}
			if !ctrlr.buffer.ShouldFlush() {
				continue
			}
			if err = ctrlr.buffer.Flush(); err != nil {
				log.Println("Error on flushing messages in buffer", err)
			}
		}
	}
}

func (ctrlr *Controller) Run() (err error) {
	go ctrlr.trackWatcherMessages()
	if err = ctrlr.watcher.Run(); err != nil {
		log.Println(err)
	}
	return
}
