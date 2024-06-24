package main

import (
	"hash/fnv"
	"sync"
	"time"
)

/*

Snowflake算法的设计目标是生成全局唯一的ID，但并不意味着它能够百分之百地避免重复。在极端情况下，如果同一毫秒内的序列号用尽，并且在下一毫秒内生成ID的速度非常快，也有可能出现重复的情况。
在你提供的代码中，当同一毫秒内的序列号用尽时，会等待下一毫秒再生成ID，以尽量避免重复。然而，如果在下一毫秒内生成ID的速度非常快，可能会导致多个goroutine在同一毫秒内生成相同的序列号，从而产生重复的ID。
为了减少重复的可能性，你可以尝试以下几个方法：
增加序列号的位数：将序列号部分的位数从12位增加到更多，例如14位或16位。这样可以增加每毫秒内可以生成的唯一ID的数量，减少序列号用尽的可能性。
控制生成ID的速度：可以通过限制生成ID的速度来降低重复的可能性。例如，可以使用time.Sleep在生成ID之间引入一定的延迟，以确保每次生成ID的时间间隔足够。
使用分布式系统：如果需要生成非常大量的ID并且要求绝对的唯一性，可以考虑使用分布式系统，其中每个节点都使用独立的Snowflake实例，并使用不同的机器ID。这样可以通过分布式的方式增加整个系统生成ID的能力，并减少重复的可能性。
需要注意的是，完全避免重复的ID是非常困难的，特别是在高并发和高速生成ID的情况下。Snowflake算法在绝大多数情况下能够提供足够的唯一性保证，但在极端情况下可能会出现重复。因此，在使用Snowflake算法生成ID时，需要根据具体的应用场景和需求来评估重复的可能性，并采取适当的措施来降低重复的风险。

*/

const (
	timestampBits  = 41
	machineIDBits  = 6
	sequenceBits   = 16
	maxMachineID   = -1 ^ (-1 << machineIDBits)
	maxSequenceID  = -1 ^ (-1 << sequenceBits)
	timeShift      = machineIDBits + sequenceBits
	machineIDShift = sequenceBits
)

// Snowflake 结构体
type Snowflake struct {
	mu         sync.Mutex
	timestamp  int64 // 时间戳部分
	machineID  int64 // 机器ID部分
	sequenceID int64 // 序列号部分
}

// NewSnowflake 创建 Snowflake 实例
func NewSnowflake(machineID int64) *Snowflake {
	return &Snowflake{
		mu:        sync.Mutex{},
		machineID: machineID,
	}
}

// GenerateID 生成全局唯一的ID
func (sf *Snowflake) GenerateID() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	currentTimestamp := time.Now().UnixNano() / 1000000

	// 如果当前时间小于等于上次生成ID的时间，需要等待下一毫秒
	if sf.timestamp >= currentTimestamp {
		time.Sleep(time.Millisecond)
		currentTimestamp = time.Now().UnixNano() / 1000000
	}

	// 如果当前时间与上次生成ID的时间相同，则需要增加序列号
	if sf.timestamp == currentTimestamp {
		sf.sequenceID = (sf.sequenceID + 1) & maxSequenceID
		if sf.sequenceID == 0 {
			// 如果序列号达到最大值，则需要等待下一毫秒
			currentTimestamp = sf.waitNextMillisecond(currentTimestamp)
		}
	} else {
		sf.sequenceID = 0
	}

	sf.timestamp = currentTimestamp

	// 生成ID
	id := (sf.timestamp << timeShift) | (sf.machineID << machineIDShift) | sf.sequenceID
	return id
}

// 等待下一毫秒
func (sf *Snowflake) waitNextMillisecond(currentTimestamp int64) int64 {
	for currentTimestamp <= sf.timestamp {
		currentTimestamp = time.Now().UnixNano() / 1000000
	}
	return currentTimestamp
}

// 将机器ID转换为int64
func machineIDToInt(machineID string) int64 {
	hash := fnv.New64a()
	hash.Write([]byte(machineID))
	hashValue := hash.Sum64()
	return int64(hashValue) & maxMachineID
}
