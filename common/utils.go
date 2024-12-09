package utils

import (
	"fmt"
	"time"

	"github.com/adamlahbib/go-ms-poc/spec"
	"golang.org/x/exp/rand"
)

func DocMsg(name string) *spec.CreateDocumentMessage {
	uid := Uid()
	doc := &spec.Document{
		Id:        id(),
		Name:      name,
		Timestamp: timestamp(),
	}
	msg := &spec.CreateDocumentMessage{
		Uid:      uid,
		Document: doc,
	}
	return msg
}

func Uid() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", t)
}

func id() string {
	return fmt.Sprintf("%d", 10000000000+rand.Intn(89999999999))
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
