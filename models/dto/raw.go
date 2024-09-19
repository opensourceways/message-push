package dto

import (
	"bytes"
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
	flatMap := make(map[string]interface{})
	flattenMap(*raw, "", flatMap)
	return flatMap
}

func flattenMap(m map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, value := range m {
		// 构造新的键名
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// 递归处理嵌套的 map
			flattenMap(v, newKey, result)
		default:
			// 直接将值添加到结果 map 中
			result[newKey] = v
		}
	}
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
	if result, ok := convertNumbers(result).(map[string]interface{}); ok {
		*raw = result
	}
}

func (flatRaw FlatRaw) ModeFilter(modeFilterJson datatypes.JSON) bool {
	if modeFilterJson == nil {
		return true
	}
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
	default:
		logrus.Errorf("convert numbers: unknown type: %T", data)
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
	parse, err := template.New("example").Parse(tmpl)
	if err != nil {
		logrus.Error(err)
	}
	t := template.Must(parse, nil)
	var resultBuffer bytes.Buffer
	err = t.Execute(&resultBuffer, raw)
	if err != nil {
		return nil
	}
	result := resultBuffer.String()
	return strings.Split(result, ",")
}
