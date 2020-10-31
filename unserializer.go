package unserializer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

const (
	ATTR_DESCRIPTION       = 1
	ATTR_EXT_FILE          = 2
	ATTR_TILE_FLAGS        = 3
	ATTR_ACTION_ID         = 4
	ATTR_UNIQUE_ID         = 5
	ATTR_TEXT              = 6
	ATTR_DESC              = 7
	ATTR_TELE_DEST         = 8
	ATTR_ITEM              = 9
	ATTR_DEPOT_ID          = 10
	ATTR_EXT_SPAWN_FILE    = 11
	ATTR_RUNE_CHARGES      = 12
	ATTR_EXT_HOUSE_FILE    = 13
	ATTR_HOUSEDOORID       = 14
	ATTR_COUNT             = 15
	ATTR_DURATION          = 16
	ATTR_DECAYING_STATE    = 17
	ATTR_WRITTENDATE       = 18
	ATTR_WRITTENBY         = 19
	ATTR_SLEEPERGUID       = 20
	ATTR_SLEEPSTART        = 21
	ATTR_CHARGES           = 22
	ATTR_CONTAINER_ITEMS   = 23
	ATTR_NAME              = 24
	ATTR_ARTICLE           = 25
	ATTR_PLURALNAME        = 26
	ATTR_WEIGHT            = 27
	ATTR_ATTACK            = 28
	ATTR_DEFENSE           = 29
	ATTR_EXTRADEFENSE      = 30
	ATTR_ARMOR             = 31
	ATTR_HITCHANCE         = 32
	ATTR_SHOOTRANGE        = 33
	ATTR_CUSTOM_ATTRIBUTES = 34
	ATTR_DECAYTO           = 35
	ATTR_WRAPID            = 36
	ATTR_STOREITEM         = 37
)

// Item defines an attribute item
type Item struct {
	Name             string                      `json:",omitempty"`
	Count            int                         `json:",omitempty"`
	Charges          int                         `json:",omitempty"`
	WrapID           int                         `json:",omitempty"`
	Text             string                      `json:",omitempty"`
	WrittenBy        string                      `json:",omitempty"`
	WrittenDate      *time.Time                  `json:",omitempty"`
	CustomAttributes map[string]*CustomAttribute `json:",omitempty"`
	Attack           int
}

// CustomAttribute defines an item attribute custom attribute
type CustomAttribute struct {
	Key   string
	Value interface{}
}

// PrettyVisualize returns a JSON visualization
func (i *Item) PrettyVisualize() string {
	v, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(v)
}

// Visualize returns a JSON visualization
func (i *Item) Visualize() string {
	v, err := json.Marshal(i)
	if err != nil {
		return err.Error()
	}
	return string(v)
}

// UnserializeHexString unserializes the given hex string
func UnserializeHexString(h string) (*Item, error) {
	data, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return Unserialize(data)
}

// Unserialize unserializes the given data
func Unserialize(data []byte) (*Item, error) {
	var ret Item

	buffer := bytes.NewBuffer(data)

	for {
		var attrType uint8
		if err := binary.Read(buffer, binary.LittleEndian, &attrType); err != nil {
			break
		}
		if attrType == 0 {
			break
		}

		switch attrType {
		case ATTR_CHARGES:
			charges, err := unserializeCharges(buffer)
			if err != nil {
				return nil, err
			}
			ret.Charges = charges
		case ATTR_ATTACK:
			attack, err := unserializeAttack(buffer)
			if err != nil {
				return nil, err
			}
			ret.Attack = attack
		case ATTR_CUSTOM_ATTRIBUTES:
			custom, err := unserializeCustomAttributes(buffer)
			if err != nil {
				return nil, err
			}
			ret.CustomAttributes = custom
		case ATTR_NAME:
			txt, err := unserializeText(buffer)
			if err != nil {
				return nil, err
			}
			ret.Name = txt
		case ATTR_WRITTENBY:
			txt, err := unserializeText(buffer)
			if err != nil {
				return nil, err
			}
			ret.WrittenBy = txt
		case ATTR_WRITTENDATE:
			d, err := unserializeDate(buffer)
			if err != nil {
				return nil, err
			}
			ret.WrittenDate = &d
		case ATTR_TEXT:
			txt, err := unserializeText(buffer)
			if err != nil {
				return nil, err
			}
			ret.Text = txt
		case ATTR_WRAPID:
			wrapId, err := unserializeWrap(buffer)
			if err != nil {
				return nil, err
			}
			ret.WrapID = wrapId
		case ATTR_COUNT, ATTR_RUNE_CHARGES:
			count, err := unserializeCount(buffer)
			if err != nil {
				return nil, err
			}
			ret.Count = count
			ret.Charges = count
		case ATTR_DEPOT_ID:
			io.CopyN(ioutil.Discard, buffer, 2)
		default:
			return nil, fmt.Errorf("Unkown attribute type: %d", attrType)
		}
	}

	return &ret, nil
}

func unserializeCustomAttribute(k string, buffer *bytes.Buffer) (*CustomAttribute, error) {
	var attrType uint8
	if err := binary.Read(buffer, binary.LittleEndian, &attrType); err != nil {
		return nil, err
	}

	var ret CustomAttribute
	ret.Key = k

	switch attrType {
	// String
	case 1:
		v, err := unserializeText(buffer)
		if err != nil {
			return nil, err
		}

		//ret.IsString = true
		ret.Value = v
	// uint64
	case 2:
		var v int64
		if err := binary.Read(buffer, binary.LittleEndian, &v); err != nil {
			return nil, err
		}

		ret.Value = v
	// double
	case 3:
		var v float64
		if err := binary.Read(buffer, binary.LittleEndian, &v); err != nil {
			return nil, err
		}

		ret.Value = v
	// bool
	case 4:
		var v bool
		if err := binary.Read(buffer, binary.LittleEndian, &v); err != nil {
			return nil, err
		}

		ret.Value = v
	}

	return &ret, nil
}

func unserializeCustomAttributes(buffer *bytes.Buffer) (map[string]*CustomAttribute, error) {
	var s uint64
	if err := binary.Read(buffer, binary.LittleEndian, &s); err != nil {
		return nil, err
	}

	list := make(map[string]*CustomAttribute, s)

	for x := uint64(0); x < s; x++ {
		attrKey, err := readString(buffer)
		if err != nil {
			return nil, err
		}

		customAttr, err := unserializeCustomAttribute(attrKey, buffer)
		if err != nil {
			return nil, err
		}

		list[attrKey] = customAttr
	}
	return list, nil
}

func unserializeAttack(buffer *bytes.Buffer) (int, error) {
	var v int32
	if err := binary.Read(buffer, binary.LittleEndian, &v); err != nil {
		return 0, err
	}
	return int(v), nil
}

func unserializeDate(buffer *bytes.Buffer) (time.Time, error) {
	var d uint32
	if err := binary.Read(buffer, binary.LittleEndian, &d); err != nil {
		return time.Time{}, err
	}

	t := time.Unix(int64(d), 0)
	return t, nil
}

func readString(buffer *bytes.Buffer) (string, error) {
	var strLen uint16
	if err := binary.Read(buffer, binary.LittleEndian, &strLen); err != nil {
		return "", err
	}

	data := make([]byte, strLen)
	if err := binary.Read(buffer, binary.LittleEndian, &data); err != nil {
		return "", err
	}

	return string(data), nil
}

func unserializeText(buffer *bytes.Buffer) (string, error) {
	return readString(buffer)
}

func unserializeCharges(buffer *bytes.Buffer) (int, error) {
	var count uint16
	if err := binary.Read(buffer, binary.LittleEndian, &count); err != nil {
		return 0, err
	}

	return int(count), nil
}

func unserializeCount(buffer *bytes.Buffer) (int, error) {
	var count uint8
	if err := binary.Read(buffer, binary.LittleEndian, &count); err != nil {
		return 0, err
	}

	return int(count), nil
}

func unserializeWrap(buffer *bytes.Buffer) (int, error) {
	var wrapId uint16
	if err := binary.Read(buffer, binary.LittleEndian, &wrapId); err != nil {
		return 0, err
	}

	return int(wrapId), nil
}
