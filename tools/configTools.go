package tools

import (
	"encoding/json"
	"net/url"

	jcs "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
	"github.com/gataca-io/vui-core/models"
)

func CopyJSON(a interface{}, b interface{}) {
	bytes, _ := json.Marshal(a)
	_ = json.Unmarshal(bytes, b)
}

func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}

func Validate(t interface{}) error {
	return models.Validate(t)
}

func CompareJSON(j1, j2 string) (bool, error) {
	j1bytes, j2bytes := []byte(j1), []byte(j2)
	jcs1, err := jcs.Transform(j1bytes)
	if err != nil {
		return false, err
	}
	jcs2, err := jcs.Transform(j2bytes)
	if err != nil {
		return false, err
	}
	return string(jcs1) == string(jcs2), nil
}

func Copy(src *url.URL) *url.URL {
	var out = new(url.URL)
	*out = *src
	return out
}

func ToMap(v interface{}) (map[string]interface{}, error) {
	mappedData := map[string]interface{}{}
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	// log.Warnf("JSON DATA: %s", bytes)
	err = json.Unmarshal(bytes, &mappedData)
	if err != nil {
		return nil, err
	}
	//log.Warnf("Unmarshaled DATA: %+v", mappedData)
	return mappedData, nil
}

func ToInterface(myMap map[string]interface{}, resultData interface{}) error {
	bytes, err := json.Marshal(myMap)
	if err != nil {
		return err
	}
	// log.Warnf("JSON DATA: %s", bytes)
	err = json.Unmarshal(bytes, resultData)
	if err != nil {
		return err
	}
	//log.Warnf("Unmarshaled DATA: %+v", mappedData)
	return nil
}
