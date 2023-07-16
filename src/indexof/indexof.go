package indexof

func Indexof[T comparable](args []T, key T) int {
	for i := 0; i < len(args); i++ {
		if args[i] == key {
			return i
		}
	}
	return -1
}
