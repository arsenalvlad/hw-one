package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	errNotTypeStruct = fmt.Errorf("entry interface not struct")
	errNoRule        = fmt.Errorf("not find rule for tag validate")
	errRuleNotMatch  = fmt.Errorf("rule for tag validate not match")
	errNoValidRule   = fmt.Errorf("rule for tag not valid")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var resultErr string
	for _, item := range v {
		resultErr += fmt.Sprintf("field: %s |||| err: %s\n", item.Field, item.Err)
	}

	return resultErr
}

func Validate(v interface{}) error {
	var errSlice ValidationErrors

	rType := reflect.TypeOf(v)
	rValue := reflect.ValueOf(v)
	if rType.Kind().String() != "struct" {
		return errNotTypeStruct
	}

	for i := 0; i < rType.NumField(); i++ {
		xType := rType.Field(i)
		xValue := rValue.Field(i)
		tagValue := xType.Tag.Get("validate")

		if tagValue == "" {
			continue
		}

		switch xType.Type.String() {
		case "string":
			err := tagStringValidate(xValue.String(), tagValue)
			if err != nil {
				errSlice = append(errSlice, ValidationError{
					Field: xType.Name,
					Err:   err,
				})
			}
		case "int":
			err := tagIntValidate(xValue.Interface().(int), tagValue)
			if err != nil {
				errSlice = append(errSlice, ValidationError{
					Field: xType.Name,
					Err:   err,
				})
			}
		case "[]int":
			for _, item := range xValue.Interface().([]int) {
				err := tagIntValidate(item, tagValue)
				if err != nil {
					errSlice = append(errSlice, ValidationError{
						Field: xType.Name,
						Err:   err,
					})
				}
			}

		case "[]string":
			for _, item := range xValue.Interface().([]string) {
				err := tagStringValidate(item, tagValue)
				if err != nil {
					errSlice = append(errSlice, ValidationError{
						Field: xType.Name,
						Err:   err,
					})
				}
			}
		default:
			continue
		}
	}

	return errSlice
}

func tagStringValidate(data string, tag string) error {
	anyRule := strings.Split(tag, "|")

	for _, value := range anyRule {
		rule := strings.Split(value, ":")
		switch rule[0] {
		case "len":
			ato, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("could not ato : %w", err)
			}

			if len(data) != ato {
				return fmt.Errorf("could not len tagStringValidate : %w", errRuleNotMatch)
			}
		case "regexp":
			matchString, err := regexp.MatchString(rule[1], data)
			if err != nil {
				return fmt.Errorf("could not rexexp match : %w", err)
			}

			if !matchString {
				return fmt.Errorf("could not regexp tagStringValidate : %w", errRuleNotMatch)
			}
		case "in":
			for _, item := range strings.Split(rule[1], ",") {
				if !strings.Contains(data, item) {
					return fmt.Errorf("could not in tagStringValidate : %w", errRuleNotMatch)
				}
			}
		default:
			return fmt.Errorf("could not tagStringValidate :%w", errNoRule)
		}
	}

	return nil
}

func tagIntValidate(data int, tag string) error {
	anyRule := strings.Split(tag, "|")

	for _, value := range anyRule {
		rule := strings.Split(value, ":")
		switch rule[0] {
		case "min":
			ato, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("could not ato : %w", err)
			}

			if data < ato {
				return fmt.Errorf("could not min tagIntValidate : %w", errRuleNotMatch)
			}
		case "max":
			ato, err := strconv.Atoi(rule[1])
			if err != nil {
				return fmt.Errorf("could not ato : %w", err)
			}

			if data > ato {
				return fmt.Errorf("could not max tagIntValidate : %w", errRuleNotMatch)
			}
		case "in":
			sliceTagValue := strings.Split(rule[1], ",")
			if len(sliceTagValue) != 2 {
				return fmt.Errorf("could not in tagIntValidate :%w", errNoValidRule)
			}

			atoMin, err := strconv.Atoi(sliceTagValue[0])
			if err != nil {
				return fmt.Errorf("could not atoMin : %w", err)
			}

			atoMax, err := strconv.Atoi(sliceTagValue[0])
			if err != nil {
				return fmt.Errorf("could not atoMax : %w", err)
			}

			if atoMin > data || atoMax < data {
				return fmt.Errorf("could not in tagIntValidate : %w", errRuleNotMatch)
			}
		default:
			return fmt.Errorf("could not tagIntValidate :%w", errNoRule)
		}
	}

	return nil
}
