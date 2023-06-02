package unionpay

import (
	"errors"
	"net/url"
	"reflect"
)

func DecodeValues(values url.Values, dst interface{}) error {
	var dstValue = reflect.ValueOf(dst)
	var dstType = dstValue.Type()
	var dstKind = dstValue.Kind()

	if dstKind == reflect.Struct {
		return errors.New("dst argument is struct")
	}

	if dstValue.IsNil() {
		return errors.New("dst argument is nil")
	}

	for {
		if dstKind == reflect.Ptr && dstValue.IsNil() {
			dstValue.Set(reflect.New(dstType.Elem()))
		}

		if dstKind == reflect.Ptr {
			dstValue = dstValue.Elem()
			dstType = dstType.Elem()
			dstKind = dstValue.Kind()
			continue
		}
		break
	}

	return decodeValues(dstType, dstValue, dstValue, values, "query")
}

func decodeValues(objType reflect.Type, parent, current reflect.Value, values url.Values, tagName string) error {
	var numField = objType.NumField()
	for i := 0; i < numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = current.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		var tag = fieldStruct.Tag.Get(tagName)
		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = fieldStruct.Name

			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			if fieldValue.Kind() == reflect.Struct {
				if err := decodeValues(fieldValue.Addr().Type().Elem(), parent, fieldValue, values, tagName); err != nil {
					return err
				}
				continue
			}
		}

		fieldValue.SetString(values.Get(tag))
	}
	return nil
}
