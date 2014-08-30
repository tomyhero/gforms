package gforms

import (
	"bytes"
	"reflect"
)

type Field interface {
	Clean(Data) (*V, error)
	Validate(*V, CleanedData) error
	Html(...RawData) string
	html(...string) string
	GetName() string
	GetWigdet() Widget
}

type ValidationError interface {
	Error() string
}

type BaseField struct {
	name       string
	validators Validators
	Widget     Widget
	Field
}

func (self *BaseField) GetName() string {
	return self.name
}

func (self *BaseField) GetWigdet() Widget {
	return self.Widget
}

func (self *BaseField) Clean(data Data) (*V, error) {
	m, hasField := data[self.GetName()]
	if hasField {
		v := m.rawValueAsString()
		m.Kind = reflect.String
		if v != nil {
			m.Value = *v
			m.IsNil = false
			return m, nil
		}
	}
	return nilV(), nil
}

func (self *BaseField) Validate(value *V, cleanedData CleanedData) error {
	if self.validators == nil {
		return nil
	}
	for _, v := range self.validators {
		err := v.Validate(value, cleanedData)
		if err != nil {
			return err
		}
	}
	return nil
}

func fieldToHtml(field Field, rds ...RawData) string {
	if len(rds) == 0 {
		if field.GetWigdet() == nil {
			return field.html()
		} else {
			return field.GetWigdet().html(field)
		}
	}
	rd := rds[0]
	v, hasField := rd[field.GetName()]
	if field.GetWigdet() == nil {
		if hasField {
			return field.html(v)
		} else {
			return field.html()
		}
	} else {
		if hasField {
			return field.GetWigdet().html(field, v)
		} else {
			return field.GetWigdet().html(field)
		}
	}
}

type templateContext struct {
	Field Field
	Value string
}

func newTemplateContext(field Field, vs ...string) templateContext {
	ctx := templateContext{
		Field: field,
	}
	if len(vs) > 0 {
		ctx.Value = vs[0]
	}
	return ctx
}

func renderTemplate(name string, ctx interface{}) string {
	var buffer bytes.Buffer
	err := Template.ExecuteTemplate(&buffer, name, ctx)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
