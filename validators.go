package gforms

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

type Validator interface {
	Name() string
	Validate(*FieldInstance, *FormInstance) error
}

type Validators []Validator

type required struct {
	Message string
	Validator
}

// Returns error if the field is not provided.
func Required(message ...string) required {
	vl := required{}
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = "この項目は必須です。"
	}
	return vl
}

func (vl required) Validate(fi *FieldInstance, fo *FormInstance) error {
	v := fi.V
	if v.IsNil || (v.Kind == reflect.String && v.Value == "") {
		return errors.New(vl.Message)
	}
	return nil
}

type maxLengthValidator struct {
	Length  int
	Message string
	Validator
}

// Returns error if the length of value is greater than length argument.
func MaxLengthValidator(length int, message ...string) maxLengthValidator {
	vl := maxLengthValidator{}
	vl.Length = length
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = fmt.Sprintf("この項目は最大 %v 文字入力して下さい。", vl.Length)
	}
	return vl
}

func (vl maxLengthValidator) Validate(fi *FieldInstance, fo *FormInstance) error {
	v := fi.V
	if v.IsNil || v.Kind != reflect.String || v.Value == "" {
		return nil
	}
	s := v.Value.(string)
	if len(s) > vl.Length {
		return errors.New(vl.Message)
	}
	return nil
}

type minLengthValidator struct {
	Length  int
	Message string
	Validator
}

// Returns error if the length of value is less than length argument.
func MinLengthValidator(length int, message ...string) minLengthValidator {
	vl := minLengthValidator{}
	vl.Length = length
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = fmt.Sprintf("この項目は最低 %v 文字以上入力して下さい。", vl.Length)
	}
	return vl
}

func (vl minLengthValidator) Validate(fi *FieldInstance, fo *FormInstance) error {
	v := fi.V
	if v.IsNil || v.Kind != reflect.String || v.Value == "" {
		return nil
	}
	s := v.Value.(string)
	if len(s) < vl.Length {
		return errors.New(vl.Message)
	}
	return nil
}

type maxValueValidator struct {
	Value   int
	Message string
	Validator
}

func MaxValueValidator(value int, message ...string) maxValueValidator {
	vl := maxValueValidator{}
	vl.Value = value
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = fmt.Sprintf("この項目は %v 以下にして下さい。", vl.Value)
	}
	return vl
}

func (vl maxValueValidator) Validate(fi *FieldInstance, fo *FormInstance) error {
	v := fi.V
	if v.IsNil || v.Kind != reflect.Int {
		return nil
	}
	iv := v.Value.(int)
	if iv > vl.Value {
		return errors.New(vl.Message)
	}
	return nil
}

type minValueValidator struct {
	Value   int
	Message string
	Validator
}

func MinValueValidator(value int, message ...string) minValueValidator {
	vl := minValueValidator{}
	vl.Value = value
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = fmt.Sprintf("この項目は %v 以上にして下さい。", vl.Value)
	}
	return vl
}

func (vl minValueValidator) Validate(fi *FieldInstance, fo *FormInstance) error {
	v := fi.V
	if v.IsNil || v.Kind != reflect.Int {
		return nil
	}
	iv := v.Value.(int)
	if iv < vl.Value {
		return errors.New(vl.Message)
	}
	return nil
}

type regexpValidator struct {
	re      *regexp.Regexp
	Message string
	Validator
}

func (vl regexpValidator) Validate(fi *FieldInstance, fo *FormInstance) error {
	v := fi.V
	if v.IsNil || v.Kind != reflect.String || v.Value == "" {
		return nil
	}
	sv := v.Value.(string)
	if !vl.re.MatchString(sv) {
		return errors.New(vl.Message)
	}
	return nil
}

// The regular expression pattern to search for the provided value.
// Returns error if regxp#MatchString is False.
func RegexpValidator(regex string, message ...string) regexpValidator {
	vl := regexpValidator{}
	re, err := regexp.Compile(regex)
	if err != nil {
		panic(err)
	}
	vl.re = re
	if len(message) > 0 {
		vl.Message = message[0]
	} else {
		vl.Message = fmt.Sprintf("正しい値を入力して下さい。")
	}
	return vl
}

// An EmailValidator that ensures a value looks like an email address.
func EmailValidator(message ...string) regexpValidator {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	if len(message) > 0 {
		return RegexpValidator(regex, message[0])
	} else {
		return RegexpValidator(regex, "正しいメールアドレスを入力して下さい。")
	}
}

// An URLValidator that ensures a value looks like an url.
func URLValidator(message ...string) regexpValidator {
	regex := `^(https?|ftp)(:\/\/[-_.!~*\'()a-zA-Z0-9;\/?:\@&=+\$,%#]+)$`
	if len(message) > 0 {
		return RegexpValidator(regex, message[0])
	} else {
		return RegexpValidator(regex, "正しいURLを入力して下さい。")
	}
}
