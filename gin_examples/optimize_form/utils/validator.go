package utils

import "strings"

// RemoveTopStruct 优化返回给用户的信息，将表单名称删除，仅保留字段名
// {'msg': {'User.password': 'password is a required field'}}
// 改为 {'msg': {'password': 'password is a required field'}}
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
