package tfsreader

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

// SerializeAsHexString serializes and returns a hex string
func (i *Item) SerializeAsHexString() (string, error) {
	v, err := i.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(v), nil
}

// Serialize serializes the given item
func (i *Item) Serialize() ([]byte, error) {
	buffer := &bytes.Buffer{}

	for _, attrType := range i.unserializedAttributes {

		if err := binary.Write(buffer, binary.LittleEndian, attrType); err != nil {
			return nil, err
		}

		switch attrType {
		case ATTR_STOREITEM:
			if err := serializeStoreItem(i, buffer); err != nil {
				return nil, err
			}
		case ATTR_CHARGES:
			if err := serializeCharges(i, buffer); err != nil {
				return nil, err
			}
		case ATTR_ATTACK:
			if err := serializeAttack(i, buffer); err != nil {
				return nil, err
			}
		case ATTR_CUSTOM_ATTRIBUTES:
			if err := serializeCustomAttributes(i, buffer); err != nil {
				return nil, err
			}
		case ATTR_NAME:
			if err := serializeText(i.Name, i, buffer); err != nil {
				return nil, err
			}
		case ATTR_WRITTENBY:
			if err := serializeText(i.WrittenBy, i, buffer); err != nil {
				return nil, err
			}
		case ATTR_WRITTENDATE:
			if err := serializeDate(i.WrittenDate, i, buffer); err != nil {
				return nil, err
			}
		case ATTR_TEXT:
			if err := serializeText(i.Text, i, buffer); err != nil {
				return nil, err
			}
		case ATTR_WRAPID:
			if err := serializeWrap(i, buffer); err != nil {
				return nil, err
			}
		case ATTR_COUNT, ATTR_RUNE_CHARGES:
			if err := serializeCount(i, buffer); err != nil {
				return nil, err
			}
		case ATTR_DEPOT_ID:
			if err := binary.Write(buffer, binary.LittleEndian, []byte{0, 0}); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("Unkown attribute type: %d", attrType)
		}
	}

	return buffer.Bytes(), nil
}

func serializeCustomAttributes(i *Item, buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.LittleEndian, uint64(len(i.CustomAttributes))); err != nil {
		return err
	}

	for _, customAttr := range i.CustomAttributes {
		if err := serializeText(customAttr.Key, i, buffer); err != nil {
			return err
		}

		if err := serializeCustomAttribute(customAttr.Value, i, buffer); err != nil {
			return err
		}
	}
	return nil
}

func serializeCustomAttribute(v interface{}, i *Item, buffer *bytes.Buffer) error {
	switch v.(type) {
	case string:
		if err := binary.Write(buffer, binary.LittleEndian, uint8(1)); err != nil {
			return err
		}
		return serializeText(v.(string), i, buffer)
	case int64:
		if err := binary.Write(buffer, binary.LittleEndian, uint8(2)); err != nil {
			return err
		}
		return binary.Write(buffer, binary.LittleEndian, v.(int64))
	case float64:
		if err := binary.Write(buffer, binary.LittleEndian, uint8(3)); err != nil {
			return err
		}
		return binary.Write(buffer, binary.LittleEndian, v.(float64))
	case bool:
		if err := binary.Write(buffer, binary.LittleEndian, uint8(4)); err != nil {
			return err
		}
		return binary.Write(buffer, binary.LittleEndian, v.(bool))
	}
	return nil
}

func serializeStoreItem(i *Item, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, i.StoreItem)
}

func serializeCount(i *Item, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, i.Count)
}

func serializeWrap(i *Item, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, i.WrapID)
}

func serializeDate(d *time.Time, i *Item, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, uint32(d.Unix()))
}

func serializeText(s string, i *Item, buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.LittleEndian, uint16(len(s))); err != nil {
		return err
	}

	if err := binary.Write(buffer, binary.LittleEndian, []byte(s)); err != nil {
		return err
	}
	return nil
}

func serializeAttack(i *Item, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, i.Attack)
}

func serializeCharges(i *Item, buffer *bytes.Buffer) error {
	return binary.Write(buffer, binary.LittleEndian, i.Charges)
}
