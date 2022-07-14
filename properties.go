package drg

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"strconv"

	"github.com/mrmarble/drg/pkg/utils"
)

type propertyType string

type Property func(io.ReadSeeker) interface{}

const (
	intProperty                     propertyType = "IntProperty"
	structProperty                  propertyType = "StructProperty"
	arrayProperty                   propertyType = "ArrayProperty"
	guidProperty                    propertyType = "Guid"
	floatProperty                   propertyType = "FloatProperty"
	dateTimeProperty                propertyType = "DateTime"
	boolProperty                    propertyType = "BoolProperty"
	multicastInlineDelegateProperty propertyType = "MulticastInlineDelegateProperty"
	setProperty                     propertyType = "SetProperty"
	mapProperty                     propertyType = "MapProperty"
	strProperty                     propertyType = "StrProperty"
	objectProperty                  propertyType = "ObjectProperty"
	uint32Property                  propertyType = "UInt32Property"
)

func properties(propertyType propertyType) Property {
	switch propertyType {
	case intProperty:
		return propertyInt
	case structProperty:
		return propertyStruct
	case arrayProperty:
		return propertyArray
	case guidProperty:
		return propertyGUID
	case floatProperty:
		return propertyFloat
	case boolProperty:
		return propertyBool
	case multicastInlineDelegateProperty:
		return propertyMulticastInlineDelegate
	case setProperty:
		return propertySet
	case mapProperty:
		return propertyMap
	case strProperty:
		return propertyStr
	case objectProperty:
		return propertyStr // This is a hack
	case uint32Property:
		return propertyInt // This is a hack

	default:
		log.Fatalf("Unsupported property type: %s", propertyType)
	}
	return nil
}

func propertyInt(r io.ReadSeeker) interface{} {
	utils.ReadNextBytes(r, 1)
	return utils.ReadNextInt32(r)
}

func propertyFloat(r io.ReadSeeker) interface{} {
	utils.ReadNextBytes(r, 1)
	return utils.ReadNextFloat32(r)
}

func propertyStruct(r io.ReadSeeker) interface{} {
	structType := utils.ReadNextString(r)
	utils.ReadNextBytes(r, 17) // Skip 16-byte empty GUID + 1-byte termination
	switch propertyType(structType) {
	case guidProperty:
		return propertyGUID(r)
	case dateTimeProperty:
		timestamp := utils.ReadNextInt64(r)

		// return time.Unix(timestamp, 0)
		return strconv.FormatInt(timestamp, 10)
	default:
		fields := make(map[string]interface{})
		for {
			if binary.LittleEndian.Uint32(utils.Peek(r, 4)) == 0 {
				break
			}
			innerName := utils.ReadNextString(r)
			if innerName == "None" {
				break
			}
			innerDataType := utils.ReadNextString(r)
			utils.ReadNextBytes(r, 8) // Skip length in int64
			property := properties(propertyType(innerDataType))
			fields[innerName] = property(r)
		}
		return fields
	}
}

func propertyGUID(r io.ReadSeeker) interface{} {
	key := utils.ReadNextBytes(r, 16)
	return hex.EncodeToString(key)
}

func propertyStructArray(r io.ReadSeeker) interface{} {
	fields := make(map[string]interface{})
	for {
		if binary.LittleEndian.Uint32(utils.Peek(r, 4)) == 0 {
			break
		}
		innerName := utils.ReadNextString(r)
		if innerName == "None" {
			break
		}
		innerDataType := utils.ReadNextString(r)
		utils.ReadNextBytes(r, 8) // Skip length in int64
		property := properties(propertyType(innerDataType))
		fields[innerName] = property(r)
	}
	return fields
}

func propertyArray(r io.ReadSeeker) interface{} {
	arrayType := utils.ReadNextString(r)
	utils.ReadNextBytes(r, 1)

	numElements := utils.ReadNextInt32(r)
	switch propertyType(arrayType) {
	case structProperty:
		utils.ReadNextString(r) // Skip array name
		dataType := utils.ReadNextString(r)
		utils.ReadNextBytes(r, 8) // Skip length in int64

		properties := []interface{}{}

		switch propertyType(dataType) {
		case structProperty:
			innerType := utils.ReadNextString(r)
			utils.ReadNextBytes(r, 17) // Skip 16-byte empty GUID + 1-byte termination
			for i := 0; i < int(numElements); i++ {
				switch propertyType(innerType) {
				case guidProperty:
					properties = append(properties, propertyGUID(r))
				default:
					properties = append(properties, propertyStructArray(r))
				}
			}
		default:
			log.Fatalf("Unsupported array type: %s", dataType)
		}

		return properties
	case intProperty:
		properties := []int32{}
		for i := 0; i < int(numElements); i++ {
			properties = append(properties, utils.ReadNextInt32(r))
		}
		return properties
	case objectProperty:
		properties := []string{}
		for i := 0; i < int(numElements); i++ {
			properties = append(properties, utils.ReadNextString(r))
		}
		return properties
	default:
		log.Fatalf("Unsupported array type: %s", arrayType)
	}
	return nil
}

func propertyBool(r io.ReadSeeker) interface{} {
	utils.ReadNextBytes(r, 1)
	return utils.ReadNextBool(r)
}

func propertyMulticastInlineDelegate(r io.ReadSeeker) interface{} {
	utils.ReadNextBytes(r, 5)
	objectPath := utils.ReadNextString(r)
	functionName := utils.ReadNextString(r)
	return struct{ ObjectPath, FunctionName string }{objectPath, functionName}
}

func propertySet(r io.ReadSeeker) interface{} {
	dataType := utils.ReadNextString(r)
	utils.ReadNextBytes(r, 5)
	numElements := utils.ReadNextInt32(r)
	properties := []interface{}{}
	for i := 0; i < int(numElements); i++ {
		switch propertyType(dataType) {
		case structProperty:
			properties = append(properties, propertyGUID(r))
		default:
			log.Fatalf("Unsupported set type: %s", dataType)

		}
	}
	return properties
}

func propertyMap(r io.ReadSeeker) interface{} {
	keyType := utils.ReadNextString(r)
	valueType := utils.ReadNextString(r)

	utils.ReadNextBytes(r, 5)
	numElements := utils.ReadNextInt32(r)
	properties := map[string]interface{}{}
	for i := 0; i < int(numElements); i++ {
		var key string
		var value interface{}

		switch propertyType(keyType) {
		case structProperty:
			key = propertyGUID(r).(string)
		case intProperty:
			key = strconv.Itoa(int(utils.ReadNextInt32(r)))
		default:
			log.Fatalf("Unsupported map key type: %s", keyType)
		}

		switch propertyType(valueType) {
		case structProperty:
			value = propertyStructArray(r)
		case intProperty:
			value = utils.ReadNextInt32(r)
		case floatProperty:
			value = utils.ReadNextFloat32(r)
		case boolProperty:
			value = utils.ReadNextBool(r)
		default:
			log.Fatalf("Unsupported map value type: %s", valueType)
		}

		properties[key] = value
	}
	return properties
}

func propertyStr(r io.ReadSeeker) interface{} {
	utils.ReadNextBytes(r, 1)
	return utils.ReadNextString(r)
}
