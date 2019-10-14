package entity

import "github.com/xeipuuv/gojsonschema"

type Schema struct {
	Location       string               `"json:location"`
	Data           string               `"json:data"`
	URN            string               `"json:urn"`
	ValidateResult *gojsonschema.Result `"json:validating-result"`
}
