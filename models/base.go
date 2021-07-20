package models

import (
	"cache2mysql/common"

	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
}

func (g *BaseModel) FileStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := common.SetField(g, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
