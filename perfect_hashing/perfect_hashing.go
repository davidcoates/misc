package misc

// Finds a perfect hash function for keys where the expected size
// of the resulting hash table is linear
func BuildHashFunction(H HashFamily, keys []interface{}) HashFunction {
	first := H.Pick(uint64(len(keys)))
	buckets := make([][]interface{}, first.Bins)
	for _, key := range keys {
		h := first.Value(key)
		buckets[h] = append(buckets[h], key)
	}
	ch := make([]chan HashFunction, first.Bins)
	for i := uint64(0); i < first.Bins; i++ {
		ch[i] = make(chan HashFunction)
		go buildInner(H, buckets[i], ch[i])
	}
	second := make([]HashFunction, first.Bins)
	bins := uint64(0)
	for i := uint64(0); i < first.Bins; i++ {
		second[i] = func(offset uint64) HashFunction {
			h := <-ch[i]
			return HashFunction{
				Bins: h.Bins,
				Value: func(key interface{}) uint64 {
					return h.Value(key) + offset
				},
			}
		}(bins)
		bins += second[i].Bins
	}
	return HashFunction{
		Bins: bins,
		Value: func(key interface{}) uint64 {
			return second[first.Value(key)].Value(key)
		},
	}
}

func buildInner(H HashFamily, keys []interface{}, out chan<- HashFunction) {
	n := uint64(len(keys))
	if n <= 1 {
		out <- HashFunction{
			Bins: n,
			Value: func(_ interface{}) uint64 {
				return 0
			},
		}
		return
	}
retry:
	h := H.Pick(n * n)
	seen := make([]bool, h.Bins)
	for _, key := range keys {
		hash := h.Value(key)
		if seen[hash] {
			goto retry
		}
		seen[hash] = true
	}
	out <- h
}
