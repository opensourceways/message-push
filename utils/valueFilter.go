package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"gorm.io/gorm/utils"
)

func ModeFilter(flatMap map[string]interface{}, modeFilterJson datatypes.JSON) bool {
	modeFilterMap := make(map[string]string)
	_ = json.Unmarshal(modeFilterJson, &modeFilterMap)
	validate := validator.New()
	for k, v := range modeFilterMap {
		err := validate.Var(flatMap[k], v)
		if err != nil {
			return false
		}
	}
	return true
}

// flattenJSON 函数将嵌套的 JSON 对象展开为平面的键值对列表
func FlattenJSON(data []byte) map[string]string {
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		fmt.Println("Error:", err)
		return make(map[string]string)
	}

	// 展开 JSON 对象
	return flattenJSON(dataMap)
}

func flattenJSON(data map[string]interface{}) map[string]string {
	flattened := make(map[string]string)
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			// 递归展开嵌套的 JSON 对象
			nested := flattenJSON(v)
			for nestedKey, nestedValue := range nested {
				flattened[key+"."+nestedKey] = nestedValue
			}
		default:
			// 将非嵌套的键值对直接添加到结果中
			flattened[key] = utils.ToString(v)
		}
	}
	return flattened
}
