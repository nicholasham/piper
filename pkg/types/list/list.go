package list

import "strings"

func Of(values ...interface{}) []interface{} {
	return values
}

func OfIntegers(values ...int) []interface{} {
	return integersToAny(values)
}

func OfStrings(values ...string) []interface{} {
	return stringsToAny(values)
}

func stringsToAny(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}

func anyToStrings(in []interface{}) []string {
	out := make([]string, len(in))
	for i := range in {
		out[i] = in[i].(string)
	}
	return out
}

func integersToAny(in []int) []interface{} {
	out := make([]interface{}, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}

func ToString(values []interface{}, separator string) string {
	return strings.Join(anyToStrings(values), separator)
}
