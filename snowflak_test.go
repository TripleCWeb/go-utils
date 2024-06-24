package utils

import (
	"fmt"
	"sync"
	"testing"
)

func TestSnowflakeGenerateID(t *testing.T) {
	snowflake := NewSnowflake(0)

	// 生成多个ID并确保它们是唯一的
	idSet := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := snowflake.GenerateID()
		if _, ok := idSet[id]; ok {
			t.Errorf("Duplicate ID generated: %d", id)
		}
		fmt.Println("id,", id)
		idSet[id] = true
	}
}

func TestSnowflakeConcurrency(t *testing.T) {
	snowflake := NewSnowflake(0)

	// 启动多个并发 goroutine 生成ID
	numGoroutines := 100
	var wg sync.WaitGroup
	idChan := make(chan int64)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				id := snowflake.GenerateID()
				fmt.Printf("i:%d j:%d, id:%d\n", idx, j, id)
				idChan <- id
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		close(idChan) // 关闭 idChan 通道
	}()

	idSet := make(map[int64]bool)
	for id := range idChan {
		if _, ok := idSet[id]; ok {
			t.Errorf("Duplicate ID generated: %d", id)
		}
		idSet[id] = true
	}

	fmt.Println(idSet)
}
