package oplog

import (
	"context"
	"log"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller test Suite")
}

var _ = Describe("Controller", func() {
	log.Println("Controller test Suite")
	var (
		newOplog *Oplog
		newCtrlr *Controller
		err      error
	)
	newOplog, _ = New()
	err = newOplog.Connect()
	newCtrlr, err = NewController(
		newOplog.srcDb,
		newOplog.srcCollections["coll_one"],
		newOplog.dstDb,
		newOplog.dstCollections["coll_one"],
	)
	newCtrlr.watcher.WatchThreshold = 2
	newCtrlr.watcher.ShouldHonorWatchThreshold = true

	go func() {
		time.Sleep(1 * time.Second)
		person := &PersonTest{Name: "Gary"}
		if _, err = newCtrlr.watcher.Collection.InsertOne(context.TODO(), person); err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Second)
		if _, err = newCtrlr.watcher.Collection.InsertOne(context.TODO(), person); err != nil {
			log.Println(err)
		}
	}()
	newCtrlr.Run()

	It("ensures no error", func() { Expect(err).To(BeNil()) })
	It("checks buffer's length", func() { Expect(newCtrlr.buffer.Length()).To(Equal(2)) })
})
