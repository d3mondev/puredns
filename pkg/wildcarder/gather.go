package wildcarder

func gatherRoots(cache *answerCache) []string {
	rootMap := make(map[string]struct{})

	for _, roots := range cache.cache {
		for _, root := range roots {
			rootMap[root] = struct{}{}
		}
	}

	found := []string{}
	for root := range rootMap {
		found = append(found, root)
	}

	return found
}
