package nbfx

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

type encoder struct {
	codec codec
}

func NewEncoder() Encoder {
	return NewEncoderWithStrings(nil)
}

func NewEncoderWithStrings(dictionaryStrings map[uint32]string) Encoder {
	encoder := &encoder{codec{make(map[uint32]string), make(map[string]uint32)}}
	if dictionaryStrings != nil {
		for k, v := range dictionaryStrings {
			encoder.codec.addDictionaryString(k, v)
		}
	}
	return encoder
}

func (e *encoder) Encode(xmlString string) ([]byte, error) {
	reader := bytes.NewReader([]byte(xmlString))
	binBuffer := &bytes.Buffer{}
	xmlDecoder := xml.NewDecoder(reader)
	token, err := xmlDecoder.RawToken()
	for err == nil {
		record := getRecordFromToken(&e.codec, token)
		if record == nil {
			return binBuffer.Bytes(), errors.New(fmt.Sprintf("Unknown Token %s", token))
		}
		err = record.write(binBuffer)
		if err != nil {
			return binBuffer.Bytes(), errors.New(fmt.Sprintf("Error writing Token %s :: %s", token, err.Error()))
		}
		token, err = xmlDecoder.RawToken()
	}
	return binBuffer.Bytes(), nil
}

func getRecordFromToken(codec *codec, token xml.Token) record {
	switch token.(type) {
	case xml.StartElement:
		return getStartElementRecordFromToken(codec, token.(xml.StartElement))
	}

	return nil
}

func getStartElementRecordFromToken(codec *codec, startElement xml.StartElement) record {
	//fmt.Printf("Getting start element for %s", startElement.Name.Local)
	prefix := startElement.Name.Space
	name := startElement.Name.Local
	prefixIndex := -1
	if len(prefix) == 1 && byte(prefix[0]) >= byte('a') && byte(prefix[0]) <= byte('z') {
		prefixIndex = int(byte(prefix[0]) - byte('a'))
	}
	var nameIndex uint32
	isNameIndexAssigned := false
	if i, ok := codec.reverseDict[name]; ok {
		nameIndex = i
		isNameIndexAssigned = true
	}

	if prefix == "" {
		if !isNameIndexAssigned {
			return &shortElementRecord{name: name}
		} else {
			return &dictionaryElementRecord{nameIndex: nameIndex}
		}
	} else if prefixIndex != -1 {
		if !isNameIndexAssigned {
			return &prefixElementAZRecord{prefixIndex: byte(prefixIndex), name: name}
		} else {
			return &prefixDictionaryElementAZRecord{prefixIndex: byte(prefixIndex), nameIndex: nameIndex}
		}
	} else {
		if !isNameIndexAssigned {
			return &elementRecord{prefix: prefix, name: name}
		} else {
			return &dictionaryElementRecord{prefix: prefix, nameIndex: nameIndex}
		}
	}
}

func writeString(reader *bytes.Reader) (string, error) {
	//var len uint32
	//strLen := len(str)
	//if err != nil {
	//	return "", err
	//}
	//strBuffer := bytes.Buffer{}
	//for i := uint32(0); i < len; {
	//	b, err := reader.ReadByte()
	//	if err != nil {
	//		return strBuffer.String(), err
	//	}
	//	strBuffer.WriteByte(b)
	//	i++
	//}
	//return strBuffer.String(), nil
	return "", nil
}
