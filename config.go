package swagger

import (
	"net/http"
	"strings"
	"reflect"
)

const (
	swaggerVersion = "2.0"
)

type SwaggerMethodParameter struct {
	Description string               `json:"description,omitempty" yaml:"description"`
	In          string               `json:"in" yaml:"in" binding:"eq=query|eq=path|eq=formData|eq=body|eq=header"`
	Name        string               `json:"name" yaml:"name" binding:"required,max=100,min=1"`
	Required    bool                 `json:"required" yaml:"required"`
	Type        string               `json:"type" yaml:"type" binding:"eq=string|eq=integer|eq=number|eq=boolean|eq=array|eq=object|eq=file"`
	Schema      *JsonSchemaObj `json:"schema,omitempty" yaml:"schema"`
}

type SwaggerMethodEntry struct {
	Description string                       `json:"description,omitempty" yaml:"description"`
	Summary     string                       `json:"summary" yaml:"summary"`
	Tags        []string                     `json:"tags" yaml:"tags" binding:"required,dive,required"`
	Parameters  []*SwaggerMethodParameter    `json:"parameters" yaml:"parameters" binding:"dive"`
	Produces    []string                     `json:"produces,omitempty" yaml:"produces"`
	Responses   map[int]*JsonSchemaObj `json:"responses" yaml:"responses" binding:"required"`
}

func SliceContain(s []string, e string) bool {
	if s == nil {
		return false
	}
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parseParameter(param *SwaggerMethodEntry, data interface{}, ptype string) {
	if reflect.TypeOf(data).Kind() != reflect.Ptr{
		panic("data must be pointer")
	}

	obj := JsonSchemaObj{}
	obj.ParseObject(data)
	if obj.Type != "object" {
		panic("form data invalid")
		return
	}
	for k, v := range obj.Properties {
		if strings.ToLower(v.Type) == "slice" ||
			strings.ToLower(v.Type) == "struct" ||
			strings.ToLower(v.Type) == "map" {
			panic("form data invalid type")
			return
		}
		parameter := &SwaggerMethodParameter{
			Description: v.Description,
			In:          ptype,
			Name:        k,
			Required:    SliceContain(obj.Required, k),
			Type:        v.Type,
		}
		param.Parameters = append(param.Parameters, parameter)
	}
}

func NewSwaggerMethodEntry(param *StructParam) *SwaggerMethodEntry {
	if param == nil {
		panic("param must exist")
		return nil
	}

	res := &SwaggerMethodEntry{}
	if len(param.Tags) < 1 {
		panic("tag must exist")
		return nil
	}
	res.Tags = param.Tags

	if len(param.Description) < 1 && len(param.Summary) < 1 {
		panic("description&summay need one at least")
		return nil
	}
	res.Summary = param.Summary
	res.Description = param.Description

	if param.ResponseData == nil {
		panic("response must exist")
		return nil
	}
	if reflect.TypeOf(param.ResponseData).Kind() != reflect.Ptr{
		panic("data must be pointer")
	}
	obj := JsonSchemaObj{}
	obj.ParseObject(param.ResponseData)
	parameter := &JsonSchemaObj{
		Schema:&obj,
	}
	resp := make(map[int]*JsonSchemaObj)
	resp[http.StatusOK] = parameter
	res.Responses = resp

	if param.JsonData != nil && param.FormData != nil {
		panic("form data and json data can not together")
		return nil
	}

	if param.JsonData != nil {
		if reflect.TypeOf(param.JsonData).Kind() != reflect.Ptr{
			panic("data must be pointer")
		}
		obj := JsonSchemaObj{}
		obj.ParseObject(param.JsonData)
		parameter := &SwaggerMethodParameter{
			Description: "Json参数",
			In:          "body",
			Name:        "body",
			Required:    true,
			Type:        "object",
			Schema:      &obj,
		}
		res.Parameters = append(res.Parameters, parameter)
	}

	if param.FormData != nil {
		parseParameter(res, param.FormData, "formData")
	}
	if param.QueryData != nil {
		parseParameter(res, param.QueryData, "query")
	}
	if param.PathData != nil {
		parseParameter(res, param.PathData, "path")
	}
	return res
}

type SwaggerDocFile map[string]SwaggerMethodEntry

type SwaggerEntry struct {
	Post   *SwaggerMethodEntry `json:"post,omitempty"`
	Get    *SwaggerMethodEntry `json:"get,omitempty"`
	Put    *SwaggerMethodEntry `json:"put,omitempty"`
	Delete *SwaggerMethodEntry `json:"delete,omitempty"`
	Patch  *SwaggerMethodEntry `json:"patch,omitempty"`
}

type SecurityDefinition struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	In          string `json:"in"`
	Name        string `json:"name"`
}

func (sentry *SwaggerEntry) SetMethod(method string, entry SwaggerMethodEntry) {
	if strings.ToLower(method) == "get" {
		sentry.Get = &entry
	} else if strings.ToLower(method) == "post" {
		sentry.Post = &entry
	} else if strings.ToLower(method) == "put" {
		sentry.Put = &entry
	} else if strings.ToLower(method) == "delete" {
		sentry.Delete = &entry
	} else if strings.ToLower(method) == "patch" {
		sentry.Patch = &entry
	} else {
		panic("invalid swagger method")
	}
}

type StructParam struct {
	FormData     interface{}
	JsonData     interface{}
	QueryData    interface{}
	PathData     interface{}
	ResponseData interface{}
	Description  string
	Summary      string
	Tags         []string
}

type SuccessResp struct {
	Message string `json:"message,omitempty" yaml:"message"`
	Result  int    `json:"result" yaml:"result"`
}
