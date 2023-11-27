package util

func JoinFields(args ...string) []string {
	l := len(args)
	ret := make([]string, 0, l/2)
	for i := 0; i < l; i += 2 {
		ret = append(ret, args[i]+"="+args[i+1])
	}
	return ret
}
