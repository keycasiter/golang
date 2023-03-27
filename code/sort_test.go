package code

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	var modelsByCondition ModelsByCondition
	models := make([]Model, 0)
	models = append(models, Model{
		Name:   "张三",
		Age:    10,
		Gender: Gender_Male,
	})
	models = append(models, Model{
		Name:   "李四",
		Age:    30,
		Gender: Gender_Male,
	})
	models = append(models, Model{
		Name:   "王五",
		Age:    30,
		Gender: Gender_Female,
	})

	modelsByCondition = append(modelsByCondition, models...)

	sort.Sort(modelsByCondition)

	for _, model := range modelsByCondition {
		fmt.Println(model)
	}
}
