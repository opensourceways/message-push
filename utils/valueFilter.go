package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
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
