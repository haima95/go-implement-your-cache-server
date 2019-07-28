package cache

type Stat struct {
	Count     int64 // 缓存目前保存的键值对数量
	KeySize   int64 // key的总字节数
	ValueSize int64 // value的总字节数
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
