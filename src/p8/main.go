package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	base := BaseMemory(client)
	fmt.Printf("base: %dB\n", base)
	for _, size := range []int64{0, 10, 20, 50, 100, 200, 1000, 5000} {
		InfoMemory(client, size, base)
	}
}

func InfoMemory(client *redis.Client, keySize int64, base int64) {
	const keyCount = 1e5
	if err := client.FlushAll().Err(); err != nil {
		panic(fmt.Sprintf("flush all db err:%v", err))
	}
	pipe := client.Pipeline()
	for i := 0; i < keyCount; i++ {
		val := make([]byte, keySize)
		for i := int64(0); i < keySize; i++ {
			val[i] = byte(rand.Intn(256))
		}
		pipe.Set(strconv.Itoa(i), string(val), 0)
	}
	if _, err := pipe.Exec(); err != nil {
		panic(fmt.Sprintf("exec pipe err: %v", err))
	}
	_ = pipe.Close()
	total := CurrentMem(client)
	total -= base
	avg := float64(total) / keyCount

	fmt.Printf("value size: %dB, total size: %dB, avg size: %.2fB\n", keySize, total, avg)
}

func BaseMemory(client *redis.Client) int64 {
	if err := client.FlushAll().Err(); err != nil {
		panic(fmt.Sprintf("flush all db err:%v", err))
	}
	return CurrentMem(client)
}

func CurrentMem(client *redis.Client) int64 {
	info := client.Info("memory").Val()
	// 输出内容的换行符为\r\n
	mem := strings.Split(info, "\n")[1]
	mem = strings.TrimSpace(mem)
	total, _ := strconv.ParseInt(strings.Split(mem, ":")[1], 10, 64)
	return total
}
