package distributed_lock

import (
	"math/rand"
)

type RedisLock struct {
	key string
	unlockValue string
}

func (r *RedisLock) New(key string) *RedisLock {
	return &RedisLock{
		key:         key,
		unlockValue: RandStr(10),
	}
}

func (r *RedisLock) Lock() {
	panic("implement me")
}

func (r *RedisLock) UnLock() {
	panic("implement me")
}

// RandStr 生成随机字符串
// length 生成字符串的长度
func RandStr(length int) string {
	data := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*()_+{}<>?,./"
	buf := make([]byte, length)
	for i:=0; i<length; i++ {
		buf[i]=data[rand.Intn(len(data))]
	}
	return string(buf)
}

