package joinmap

// JoinMap joins two map
func JoinMap(map2, map1 map[string]string) map[string]string {
	for key, value := range map1 {
		_, ok := map2[key]
		if ok {
			map2["_"+key] = value
		} else {
			map2[key] = value
		}
	}
	return map2
}
