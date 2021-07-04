package utils

func DeleteArrayElement(arr []string, elem string) []string {
	res := []string{}
	for i, v := range arr {
		if v == elem {
			res = append(res, arr[:i]...)
			res = append(res, arr[i+1:]...)
			return res
		}
	}
	return arr
}
