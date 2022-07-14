package drg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/mrmarble/drg/pkg/utils"
)

const (
	Header = "GVAS"
)

var (
	ErrInvalidHeader = errors.New("Invalid header")
	ErrInvalidOffset = errors.New("Invalid offset")
)

type EngineVersion struct {
	Major   uint16
	Minor   uint16
	Patch   uint16
	Build   uint32
	BuildID string
}

type Metadata struct {
	SaveVersion         int32
	PackageVersion      int32
	EngineVersion       EngineVersion
	CustomFormatVersion int32
	CustomFormatData    map[string]int32
	SaveGameType        string
}

func DecodeMetadata(r io.Reader) (*Metadata, error) {
	header := utils.ReadNextBytes(r, len(Header))
	if string(header) != Header {
		return nil, ErrInvalidHeader
	}
	var metadata Metadata

	err := utils.ReadStruct(r, &metadata)
	if err != nil {
		return nil, fmt.Errorf("Failed to read metadata: %v", err)
	}

	return &metadata, nil
}

func Decode(r io.ReadSeeker) (map[string]interface{}, error) {
	offset, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, ErrInvalidOffset
	}

	if offset == 0 {
		_, err := DecodeMetadata(r)
		if err != nil {
			return nil, err
		}
	}

	return decode(r), nil
}

func decode(r io.ReadSeeker) map[string]interface{} {
	fields := make(map[string]interface{})
	for {
		if binary.LittleEndian.Uint32(utils.Peek(r, 4)) == 0 {
			break
		}
		name := utils.ReadNextString(r)
		if name == "None" {
			break
		}
		dataType := utils.ReadNextString(r)
		utils.ReadNextBytes(r, 8) // Skip length in int64

		property := properties(propertyType(dataType))
		fields[name] = property(r)
	}
	return fields
}
