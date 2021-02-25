package protobuf

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

// ErrNotFound is returned when the requested message type is not found
// in the protobuf files.
var ErrNotFound = errors.New("not found")

// ProtoBuf provides facilities for dynamically marshaling/unmarshalling
// protobuf messages.
type ProtoBuf struct {
	Files       []string `json:"files"`
	ImportPaths []string `json:"import_paths"`
	MessageType string   `json:"message_type"`

	once sync.Once
	md   *desc.MessageDescriptor
}

// Init initialises the protobuf parser. Only first invocation is processed.
func (pb *ProtoBuf) Init() error {
	var err error
	pb.once.Do(func() {
		pb.md, err = protobufParse(pb.MessageType, pb.Files, pb.ImportPaths)
	})
	return err
}

// Unmarshal reads and parses the protobuf serialised data into an instance
// of protobuf message.
func (pb *ProtoBuf) Unmarshal(d []byte) (*dynamic.Message, error) {
	if err := pb.Init(); err != nil {
		return nil, err
	} else if pb.md == nil {
		return nil, fmt.Errorf("%w: '%s'", ErrNotFound, pb.MessageType)
	}

	msg := dynamic.NewMessage(pb.md)
	if err := msg.Unmarshal(d); err != nil {
		return nil, err
	}
	return msg, nil
}

func protobufParse(msgType string, files, importDirs []string) (*desc.MessageDescriptor, error) {
	p := &protoparse.Parser{
		InferImportPaths: true,
		ImportPaths:      importDirs,
	}

	descriptors, err := p.ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %w", err)
	}

	for _, fd := range descriptors {
		msgDesc := fd.FindMessage(msgType)
		if msgDesc != nil {
			return msgDesc, nil
		}
	}

	return nil, nil
}
