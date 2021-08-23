package mackerel

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandStringListFromSet(set *schema.Set) []string {
	strings := make([]string, 0, set.Len())
	for _, v := range set.List() {
		strings = append(strings, v.(string))
	}
	return strings
}

func flattenStringList(strings []string) []interface{} {
	vs := make([]interface{}, 0, len(strings))
	for _, v := range strings {
		vs = append(vs, v)
	}
	return vs
}

func flattenStringListToSet(strings []string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(strings))
}
