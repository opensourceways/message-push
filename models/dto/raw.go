package dto

import (
	"bytes"
	flattener "github.com/anshal21/json-flattener"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/opensourceways/message-push/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"strings"
	"text/template"
)

type FlatRaw map[string]interface{}

type Raw map[string]interface{}

type FlatRawString map[string]string

func (raw *Raw) Flatten() FlatRaw {
	s, _ := json.Marshal(raw)
	flatJSON, _ := flattener.FlattenJSON(string(s), flattener.DotSeparator)
	flatMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(flatJSON), &flatMap)
	return flatMap
}

func (raw *Raw) FromJson(jsonStr []byte) {
	// 创建一个Decoder对象，并调用UseNumber()
	decoder := json.NewDecoder(bytes.NewReader(jsonStr))
	decoder.UseNumber()

	// 解码到map
	var result map[string]interface{}
	err := decoder.Decode(&result)
	if err != nil {
		logrus.Error(err)
		return
	}
	result = convertNumbers(result).(map[string]interface{})
	*raw = result
}

func (flatRaw FlatRaw) ModeFilter(modeFilterJson datatypes.JSON) bool {
	modeFilterMap := make(map[string]string)
	_ = json.Unmarshal(modeFilterJson, &modeFilterMap)
	validate := validator.New()
	for k, v := range modeFilterMap {
		err := validate.Var(flatRaw[k], v)
		if err != nil {
			return false
		}
	}
	return true
}

func convertNumbers(data interface{}) interface{} {
	switch data := data.(type) {
	case json.Number:
		if i, err := data.Int64(); err == nil {
			return i
		} else if f, err := data.Float64(); err == nil {
			return f
		}
	case map[string]interface{}:
		for k, v := range data {
			data[k] = convertNumbers(v)
		}
	case []interface{}:
		for i, v := range data {
			data[i] = convertNumbers(v)
		}
	}
	return data
}

func (flatRaw *FlatRaw) StringifyMap() map[string]string {
	result := make(map[string]string)
	for key, value := range *flatRaw {
		result[key] = utils.ToString(value)
	}
	return result
}

func (raw *Raw) ToMessageArgs(messageTemplate string) []string {
	tmpl := messageTemplate
	t := template.Must(template.New("example").Parse(tmpl))
	var resultBuffer bytes.Buffer
	err := t.Execute(&resultBuffer, raw)
	if err != nil {
		return nil
	}
	result := resultBuffer.String()
	return strings.Split(result, ",")
}
