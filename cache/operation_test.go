package cache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type Value struct {
	Name string
	Age  int
}

func Test_cache_Get(t *testing.T) {
	c := New()
	v := Value{"库陈胜", 2022}
	c.Set("test", v)
	c.Set("test1", "啊哈")
	get, err := c.Get("test")
	t.Logf("result %v,%v", get, err)

	c.Remove("test")
	get, err = c.Get("test")
	t.Logf("result %v,%v", get, err)
}

func Test_cache_Get1(t *testing.T) {
	c := NewWithExpiration(1 * time.Second)
	c.Set("testV", "库陈胜")
	time.Sleep(1 * time.Second)
	get, b := c.Get("testV")
	if b {
		t.Error("expire fail get is ", get)
	}
	c.SetHash("Key", "subKey", "库陈胜")
	c.SetHash("Key", "6不liuKey@", "贼溜")

	if v1, b := c.Get("Key"); !b {
		t.Error("未获取到hash值")
	} else {
		t.Logf("获取到hash值 %v", v1)
	}
	time.Sleep(1 * time.Second)

	if v2, b := c.GetHash("Key", "subKey"); !b {
		t.Error("未获取到hash值 v2")
	} else {
		t.Logf("获取到hash 及子值 %v", v2)
	}
	time.Sleep(2 * time.Second)

	if _, b = c.Get("Key"); b {
		t.Error("Key 未过期")
	}

	if _, b = c.GetHash("Key", "subKey"); b {
		t.Error("Key - subKey 未过期")
	}
}

//性能测试
//fixme 并发问题有待处理
func Test_cache_Get2(t *testing.T) {
	c := NewWithExpiration(2 * time.Second)
	ch := make(chan int8, 20000)
	start := time.Now()
	println("开始执行", start.UnixMilli())
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("%s%d", "Key", i)
		go func() {
			c.Set(key, "库陈胜"+key)
			ch <- int8(1)
			c.SetHash(key+"hash", strconv.Itoa(i), "性能测试"+key)
			ch <- int8(1)
		}()
	}
	println("PUT结束执行,耗时", time.Now().UnixMilli()-start.UnixMilli(), "ms", "key总数", c.Cap())
	ch1 := make(chan int8, 20000)
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("%s%d", "Key", i)
		go func() {
			if v, b := c.Get(key); b {
				println("key=", key, "value = ", v.(string))
			}
			ch1 <- int8(1)
			subKey := strconv.Itoa(i)
			if v, b := c.GetHash(key+"hash", subKey); b {
				println("key=", key, "subKey=", subKey, "value = ", v.(string))
			}
			ch1 <- int8(1)
		}()
	}

	println("当前有多少key?", c.Cap())
	times := 1
	for c.Cap() > 0 {
		time.Sleep(1 * time.Second)
		println("沉睡", times, "秒后，剩余多少Key?", c.Cap())
		times++
	}
}
