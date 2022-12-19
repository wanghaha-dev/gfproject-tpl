package utils

import (
	"github.com/wanghaha-dev/gf/frame/g"
	"strings"
)

// GetAddFieldsRule 获取表字段的校验规则
func GetAddFieldsRule(table string) (map[string]string, error) {
	fields, err := g.Model(table).TableFields(table)
	if err != nil {
		return nil, err
	}

	mapRules := make(map[string]string)
	for key, val := range fields {
		// 不需要验证的字段跳过
		if key == "created_at" || key == "updated_at" || key == "deleted_at" {
			continue
		}

		var rules []string
		if !val.Null {
			rules = append(rules, "required")
		}
		if strings.Contains(val.Type, "int") {
			rules = append(rules, "integer")
		}
		// 组合校验规则
		mapRules[key] = strings.Join(rules, "|")
	}
	return mapRules, nil
}

// GetUpdateFieldsRule 获取表字段的校验规则
func GetUpdateFieldsRule(table string) (map[string]string, error) {
	fields, err := g.Model(table).TableFields(table)
	if err != nil {
		return nil, err
	}

	mapRules := make(map[string]string)
	for key, val := range fields {
		// 不需要验证的字段跳过
		if key == "id" || key == "created_at" || key == "updated_at" || key == "deleted_at" {
			continue
		}

		var rules []string
		if !val.Null {
			rules = append(rules, "required")
		}
		if strings.Contains(val.Type, "int") {
			rules = append(rules, "integer")
		}
		// 组合校验规则
		mapRules[key] = strings.Join(rules, "|")
	}
	return mapRules, nil
}
