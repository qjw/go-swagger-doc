package swagger

import (
	"reflect"
	"strings"
)

type JsonSchemaObj struct {
	Description string                    `json:"description,omitempty" yaml:"description"`
	Type        string                    `json:"type,omitempty" yaml:"type" binding:"required,min=1"`
	Items       *JsonSchemaObj            `json:"items,omitempty" yaml:"items"`
	Properties  map[string]*JsonSchemaObj `json:"properties,omitempty" yaml:"properties"`
	Required    []string                  `json:"required,omitempty" yaml:"required"`
	Schema      *JsonSchemaObj            `json:"schema,omitempty" yaml:"schema"`
}

func (obj *JsonSchemaObj) ParseObject(variable interface{}) {
	value := reflect.ValueOf(variable)
	obj.read(value.Type(), "")
}

var formatMapping = map[string][]string{
	"time.Time": []string{"string", "date-time"},
}

var kindMapping = map[reflect.Kind]string{
	reflect.Bool:    "boolean",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Uint:    "integer",
	reflect.Uint8:   "integer",
	reflect.Uint16:  "integer",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.String:  "string",
	reflect.Slice:   "array",
	reflect.Struct:  "object",
	reflect.Map:     "object",
}

func getTypeFromMapping(t reflect.Type) (string, string, reflect.Kind) {
	//if v, ok := formatMapping[t.String()]; ok {
	//	return v[0], v[1], reflect.String
	//}

	if v, ok := kindMapping[t.Kind()]; ok {
		return v, "", t.Kind()
	}

	return "", "", t.Kind()
}

func (obj *JsonSchemaObj) read(t reflect.Type, doc string) {
	jsType, _, kind := getTypeFromMapping(t)
	if jsType != "" {
		obj.Type = jsType
	}
	obj.Description = doc
	//if format != "" {
	//	obj.Format = format
	//}

	switch kind {
	case reflect.Slice:
		obj.readFromSlice(t)
	case reflect.Map:
		obj.readFromMap(t)
	case reflect.Struct:
		obj.readFromStruct(t)
	case reflect.Ptr:
		obj.read(t.Elem(), doc)
	}
}

func (obj *JsonSchemaObj) readFromSlice(t reflect.Type) {
	jsType, _, kind := getTypeFromMapping(t.Elem())
	if kind == reflect.Uint8 {
		obj.Type = "string"
	} else if jsType != "" {
		obj.Items = &JsonSchemaObj{Type: jsType}
		obj.Items.read(t.Elem(), "")
	}
}

func (obj *JsonSchemaObj) readFromMap(t reflect.Type) {
	jsType, _, _ := getTypeFromMapping(t.Elem())

	if jsType != "" {
		obj.Properties = make(map[string]*JsonSchemaObj, 0)
		var tmp_obj = &JsonSchemaObj{Type: jsType}
		obj.Properties[".*"] = tmp_obj
		tmp_obj.read(t.Elem(), "")
	}
}

func (obj *JsonSchemaObj) readFromStruct(t reflect.Type) {
	obj.Type = "object"
	obj.Properties = make(map[string]*JsonSchemaObj, 0)

	count := t.NumField()
	for i := 0; i < count; i++ {
		field := t.Field(i)
		if field.Anonymous {
			obj.read(field.Type,"")
			continue
		}

		tag := field.Tag.Get("json")
		name, opts := parseTag(tag)
		if name == "" {
			name = field.Name
		}
		if name == "-" {
			continue
		}

		var tmp_obj = &JsonSchemaObj{}
		obj.Properties[name] = tmp_obj
		tmp_obj.read(field.Type, field.Tag.Get("doc"))

		if !opts.Contains("omitempty") {
			obj.Required = append(obj.Required, name)
		}
	}
}

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

type tagOptions string
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}

	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func (o tagOptions) GetValue(optionName string) (string, bool) {
	if len(o) == 0 {
		return "", false
	}

	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if strings.ContainsAny(s, "=") {
			var strs = strings.Split(s, "=")
			if len(strs) == 2 && strs[0] == optionName {
				return strs[1], true
			}
		}
		s = next
	}
	return "", false
}
