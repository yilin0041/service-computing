package myjson

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"sort"
	"strconv"
)

type marshalData struct {
	bytes.Buffer
}

//Marshal :the most important function
func Marshal(v interface{}) ([]byte, error) {
	var data marshalData
	val := reflect.ValueOf(v)
	if err := data.marshal(val); err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func (data *marshalData) marshal(val reflect.Value) error {
	if val.Kind() == reflect.Int {
		return data.marshalInt(val)
	} else if val.Kind() == reflect.String {
		return data.marshalString(val)
	} else if val.Kind() == reflect.Slice {
		return data.marshalSlice(val)
	} else if val.Kind() == reflect.Array {
		return data.marshalArray(val)
	} else if val.Kind() == reflect.Map {
		return data.marshalMap(val)
	} else if val.Kind() == reflect.Struct {
		return data.marshalStruct(val)
	} else if val.Kind() == reflect.Ptr {
		return data.marshalPtr(val)
	} else if val.Kind() == reflect.Interface {
		return data.marshalInterface(val)
	} else {
		return errors.New("Unknown type  " + val.Kind().String())
	}
}

func (data *marshalData) marshalInt(val reflect.Value) error {
	if _, err := data.Write(strconv.AppendInt([]byte{}, val.Int(), 10)); err != nil {
		return err
	}
	return nil
}

func (data *marshalData) marshalString(val reflect.Value) error {
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	if _, err := data.Write([]byte(val.String())); err != nil {
		return err
	}
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	return nil
}

func (data *marshalData) marshalSlice(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		if err := data.WriteByte('['); err != nil {
			return err
		}
		for i := 0; i < val.Len(); i++ {
			element := reflect.ValueOf(val.Index(i).Interface())
			if err := data.marshal(element); err != nil {
				return err
			}
			if i != val.Len()-1 {
				if err := data.WriteByte(','); err != nil {
					return err
				}
			}
		}
		return data.WriteByte(']')
	}
	valBytes := val.Bytes()
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	if _, err := data.Write(valBytes); err != nil {
		return err
	}
	if err := data.WriteByte('"'); err != nil {
		return err
	}
	if err := data.WriteByte(':'); err != nil {
		return err
	}
	return nil
}

func (data *marshalData) marshalArray(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		if err := data.WriteByte('['); err != nil {
			return err
		}
		for i := 0; i < val.Len(); i++ {
			element := reflect.ValueOf(val.Index(i).Interface())
			if err := data.marshal(element); err != nil {
				return err
			}
			if i != val.Len()-1 {
				if err := data.WriteByte(','); err != nil {
					return err
				}
			}
		}
		return data.WriteByte(']')
	}
	return errors.New("Unknown type  " + elemKind.String())
}

type sortableByteSliceSlice [][]byte

func (s sortableByteSliceSlice) Len() int {
	return len(s)
}

func (s sortableByteSliceSlice) Less(i, j int) bool {
	return bytes.Compare(s[i], s[j]) == -1
}

func (s sortableByteSliceSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func (data *marshalData) marshalMap(val reflect.Value) error {
	keys := val.MapKeys()
	if err := data.WriteByte('{'); err != nil {
		return err
	}
	raw := make(sortableByteSliceSlice, len(keys))
	for i, key := range keys {
		if key.Kind() != reflect.String {
			e := "Map keys must be 'string' type,your keys is '" + (key.Kind().String()) + "' type"
			return errors.New(e)
		}
		raw[i] = []byte(key.String())
	}
	sort.Sort(raw)
	i := 0
	for _, rawKey := range raw {
		key := string(rawKey)
		vKey := reflect.ValueOf(key)
		if err := data.marshal(vKey); err != nil {
			return err
		}
		if err := data.WriteByte(':'); err != nil {
			return err
		}
		value := val.MapIndex(vKey)
		if err := data.marshal(value); err != nil {
			return err
		}
		if i != raw.Len()-1 {
			if err := data.WriteByte(','); err != nil {
				return err
			}
		}
		i++
	}
	return data.WriteByte('}')
}

type positionedField struct {
	name []byte
	pos  int
}

type positionedFieldsByName []positionedField

func myTag(value reflect.Value, name string) string {
	var tag string
	field, hasField := value.Type().FieldByName(name)
	if !hasField {
		tag = ""
	} else {
		tag = string(field.Tag)
	}
	var readTag *string
	const fieldRegexp = `json:"([\w- ]*)"`
	reg := regexp.MustCompile(fieldRegexp)
	if matches := reg.FindStringSubmatch(tag); len(matches) > 2 {
		panic("regexp returns more then two groups!")
	} else if len(matches) == 2 {
		readTag = &matches[1]
	} else {
		readTag = nil
	}
	if readTag == nil {
		return name
	} else if *readTag == "" || *readTag == "-" {
		return ""
	} else {
		return *readTag
	}
}
func (data *marshalData) marshalStruct(val reflect.Value) error {
	if err := data.WriteByte('{'); err != nil {
		return err
	}
	valType := val.Type()

	fields := positionedFieldsByName{}
	count := 0
	for i := 0; i < val.NumField(); i++ {
		fieldOpt := myTag(val, valType.Field(i).Name)
		if len(fieldOpt) == 0 {
			count++
			continue
		}
		temp := (string)(valType.Field(i).Name)
		if temp[0] < 'A' || temp[0] > 'Z' {
			count++
			continue
		}
		fields = append(fields, positionedField{[]byte(fieldOpt), i})
	}
	for _, f := range fields {
		count++
		if err := data.marshal(reflect.ValueOf(f.name)); err != nil {
			return err
		}
		if err := data.marshal(val.Field(f.pos)); err != nil {
			return err
		}
		if count < val.NumField() {
			if err := data.WriteByte(','); err != nil {
				return err
			}
		}
	}
	return data.WriteByte('}')
}

func (data *marshalData) marshalPtr(val reflect.Value) error {
	return data.marshal(val.Elem())
}

func (data *marshalData) marshalInterface(val reflect.Value) error {
	return data.marshal(val.Elem())
}
