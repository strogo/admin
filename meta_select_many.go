package admin

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"

	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
)

// CompositePrimaryKey the string that represents the composite primary key
const CompositePrimaryKey = "CompositePrimaryKey"

// CompositePrimaryKeyField to embed into the struct that requires composite primary key in select many
type CompositePrimaryKeyField struct {
	CompositePrimaryKey string `gorm:"-"`
}

// SetCompositePrimaryKey set CompositePrimaryKey in a specific format
func (cpk *CompositePrimaryKeyField) SetCompositePrimaryKey(id uint, versionName string) {
	cpk.CompositePrimaryKey = GenCompositePrimaryKey(id, versionName)
}

func GenCompositePrimaryKey(id uint, versionName string) string {
	return fmt.Sprintf("%d%s%s", id, resource.CompositePrimaryKeySeparator, versionName)
}

// SelectManyConfig meta configuration used for select many
type SelectManyConfig struct {
	Collection               interface{} // []string, [][]string, func(interface{}, *qor.Context) [][]string, func(interface{}, *admin.Context) [][]string
	DefaultCreating          bool
	Placeholder              string
	SelectionTemplate        string
	SelectMode               string // select, select_async, bottom_sheet
	Select2ResultTemplate    template.JS
	Select2SelectionTemplate template.JS
	ForSerializedObject      bool
	RemoteDataResource       *Resource
	RemoteDataHasImage       bool
	PrimaryField             string
	SelectOneConfig
}

// GetTemplate get template for selection template
func (selectManyConfig SelectManyConfig) GetTemplate(context *Context, metaType string) ([]byte, error) {
	if metaType == "form" && selectManyConfig.SelectionTemplate != "" {
		return context.Asset(selectManyConfig.SelectionTemplate)
	}
	return nil, errors.New("not implemented")
}

// ConfigureQorMeta configure select many meta
func (selectManyConfig *SelectManyConfig) ConfigureQorMeta(metaor resource.Metaor) {
	if meta, ok := metaor.(*Meta); ok {
		selectManyConfig.SelectOneConfig.Collection = selectManyConfig.Collection
		selectManyConfig.SelectOneConfig.SelectMode = selectManyConfig.SelectMode
		selectManyConfig.SelectOneConfig.DefaultCreating = selectManyConfig.DefaultCreating
		selectManyConfig.SelectOneConfig.Placeholder = selectManyConfig.Placeholder
		selectManyConfig.SelectOneConfig.RemoteDataResource = selectManyConfig.RemoteDataResource
		selectManyConfig.SelectOneConfig.PrimaryField = selectManyConfig.PrimaryField

		selectManyConfig.SelectOneConfig.ConfigureQorMeta(meta)

		selectManyConfig.RemoteDataResource = selectManyConfig.SelectOneConfig.RemoteDataResource
		selectManyConfig.SelectMode = selectManyConfig.SelectOneConfig.SelectMode
		selectManyConfig.DefaultCreating = selectManyConfig.SelectOneConfig.DefaultCreating
		selectManyConfig.PrimaryField = selectManyConfig.SelectOneConfig.PrimaryField
		meta.Type = "select_many"

		// Set FormattedValuer
		if meta.FormattedValuer == nil {
			meta.SetFormattedValuer(func(record interface{}, context *qor.Context) interface{} {
				reflectValues := reflect.Indirect(reflect.ValueOf(meta.GetValuer()(record, context)))
				var results []string
				if reflectValues.IsValid() {
					for i := 0; i < reflectValues.Len(); i++ {
						results = append(results, utils.Stringify(reflectValues.Index(i).Interface()))
					}
				}
				return results
			})
		}
	}
}
