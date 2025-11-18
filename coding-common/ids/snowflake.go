package ids

import (
	"fmt"
	"sync"
	"time"
)

const (
	twepoch          int64 = 1420041600000
	workerIdBits           = 5
	datacenterIdBits       = 5
	maxWorkerId            = -1 ^ (-1 << workerIdBits)
	maxDatacenterId        = -1 ^ (-1 << datacenterIdBits)
	sequenceBits     uint8 = 12

	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
)

type SnowflakeIdWorker struct {
	mu            sync.Mutex
	workerId      int64
	datacenterId  int64
	sequence      int64
	lastTimestamp int64
}

func NewSnowflakeIdWorker(workerId int64, datacenterId int64) (*SnowflakeIdWorker, error) {
	if workerId > maxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker Id can't be greater than %d or less than 0", maxWorkerId)
	}
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenter Id can't be greater than %d or less than 0", maxDatacenterId)
	}

	return &SnowflakeIdWorker{
		workerId:      workerId,
		datacenterId:  datacenterId,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

func (w *SnowflakeIdWorker) NextId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	timestamp := timeGen()
	lastTimestamp := w.lastTimestamp
	sequence := w.sequence
	// 如果当前时间小于上一次ID生成的时间戳, 说明系统时钟回退过这个时候应当抛出异常
	if timestamp < lastTimestamp {
		panic(
			fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds", lastTimestamp-timestamp),
		)
	}

	// 如果是同一时间生成的, 则进行毫秒内序列
	if lastTimestamp == timestamp {
		sequence = (sequence + 1) & sequenceMask
		// 毫秒内序列溢出
		if sequence == 0 {
			// 阻塞到下一个毫秒,获得新的时间戳
			timestamp = tilNextMillis(lastTimestamp)
		}
	} else {
		// 时间戳改变, 毫秒内序列重置
		sequence = 0
	}

	// 上次生成ID的时间戳
	lastTimestamp = timestamp

	// 移位并通过或运算拼到一起组成64位的ID
	return ((timestamp - twepoch) << timestampLeftShift) | (w.datacenterId << datacenterIdShift) | (w.workerId << workerIdShift) | sequence
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}
func timeGen() int64 {
	return time.Now().UnixMilli()
}
