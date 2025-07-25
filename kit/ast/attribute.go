package ast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tomdoesdev/knox/kit/errs"
)

var (
	ErrNoAttributeValue    = errs.New(NoValueErrCode, "no attribute value")
	ErrTypeAssertionFailed = errs.New(InvalidTypeCast, "failed to assert type")
)

type Attribute struct {
	value any
}

func NewAttribute(value any) *Attribute {
	return &Attribute{value}
}

func (a *Attribute) AsBool() (bool, error) {
	if a.value == nil {
		return false, ErrNoAttributeValue
	}

	switch v := a.value.(type) {
	case bool:
		return v, nil
	case *bool:
		return *v, nil
	case int:
		return v != 0, nil
	case *int:
		return *v != 0, nil
	case int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case *int8, *int16, *int32, *int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case string:
		s, err := strconv.ParseBool(v)
		if err != nil {
			return false, errs.Wrap(err, InvalidTypeCast, "failed to convert string to bool")
		}
		return s, nil
	case *string:
		s, err := strconv.ParseBool(*v)
		if err != nil {
			return false, errs.Wrap(err, InvalidTypeCast, "failed to convert string to bool")
		}
		return s, nil
	default:
		return false, ErrTypeAssertionFailed.WithContext("as", "bool")
	}
}

func (a *Attribute) AsString() (string, error) {
	if a.value == nil {
		return "", ErrNoAttributeValue
	}

	switch v := a.value.(type) {
	case string:
		return v, nil
	case *string:
		return *v, nil
	case fmt.Stringer:
		return v.String(), nil
	case bool:
		return strconv.FormatBool(v), nil
	case *bool:
		return strconv.FormatBool(*v), nil
	case int:
		return strconv.Itoa(v), nil
	case *int:
		return strconv.Itoa(*v), nil
	case int8, int16, int32, int64:
		return fmt.Sprintf("%d", reflect.ValueOf(v).Int()), nil
	case *int8, *int16, *int32, *int64:
		return fmt.Sprintf("%d", reflect.ValueOf(v).Elem().Int()), nil
	case float32, float64:
		return fmt.Sprintf("%g", reflect.ValueOf(v).Float()), nil
	case *float32, *float64:
		return fmt.Sprintf("%g", reflect.ValueOf(v).Elem().Float()), nil
	default:
		return "", ErrTypeAssertionFailed.WithContext("as", "string")
	}
}

func (a *Attribute) AsInt() (int, error) {
	if a.value == nil {
		return 0, ErrNoAttributeValue
	}

	switch v := a.value.(type) {
	case int:
		return v, nil
	case *int:
		return *v, nil
	case int8, int16, int32, int64:
		return int(reflect.ValueOf(v).Int()), nil
	case *int8, *int16, *int32, *int64:
		return int(reflect.ValueOf(v).Elem().Int()), nil
	case float32, float64:
		return int(reflect.ValueOf(v).Float()), nil
	case *float32, *float64:
		return int(reflect.ValueOf(v).Elem().Float()), nil
	case string:
		result, err := strconv.Atoi(v)
		if err != nil {
			return 0, errs.Wrap(err, InvalidTypeCast, "failed to convert string to int")
		}
		return result, nil
	case *string:
		result, err := strconv.Atoi(*v)
		if err != nil {
			return 0, errs.Wrap(err, InvalidTypeCast, "failed to convert string to int")
		}
		return result, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case *bool:
		if *v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ErrTypeAssertionFailed.WithContext("as", "int")
	}
}

func (a *Attribute) AsFloat() (float64, error) {
	if a.value == nil {
		return 0.0, ErrNoAttributeValue
	}

	switch v := a.value.(type) {
	case float64:
		return v, nil
	case *float64:
		return *v, nil
	case float32:
		return float64(v), nil
	case *float32:
		return float64(*v), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int()), nil
	case *int, *int8, *int16, *int32, *int64:
		return float64(reflect.ValueOf(v).Elem().Int()), nil
	case string:
		result, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0.0, errs.Wrap(err, InvalidTypeCast, "failed to convert string to float")
		}
		return result, nil
	case *string:
		result, err := strconv.ParseFloat(*v, 64)
		if err != nil {
			return 0.0, errs.Wrap(err, InvalidTypeCast, "failed to convert string to float")
		}
		return result, nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	case *bool:
		if *v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0.0, ErrTypeAssertionFailed.WithContext("as", "float64")
	}
}
