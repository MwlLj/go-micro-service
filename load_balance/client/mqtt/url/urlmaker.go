package url

import (
	"bytes"
)

type IUrlMaker interface {
	Make(topic *string) *string
}

type CUrlMaker struct {
	ServerUniqueCode *string
}

func (this *CUrlMaker) Make(topic *string) *string {
	if this.ServerUniqueCode == nil {
		return topic
	}
	// if *this.ServerUniqueCode == "" -> join / must
	var buffer bytes.Buffer
	buffer.WriteString(*this.ServerUniqueCode)
	bTopic := []byte(*topic)
	if bTopic[0] != '/' {
		buffer.WriteString("/")
	}
	buffer.WriteString(*topic)
	s := buffer.String()
	return &s
}
