package gid

import (
	"sync"
	"time"
)

/* ================================================================================
 * snowflake id
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

type (
	Snowflake struct {
		startTimestamp int64
		nodeBits       uint8
		seqBits        uint8
		nodeMax        int64
		nodeMask       int64
		timeShift      uint8
		nodeShift      uint8
		seqMask        int64
	}

	SnowflakeNode struct {
		timestamp    int64
		workerNodeId int64
		seq          int64
		mu           sync.Mutex
		snowflake    *Snowflake
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化Snowflake
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewSnowflake() *Snowflake {
	snowflake := &Snowflake{
		startTimestamp: 1551377400197,
		nodeBits:       10,
		seqBits:        12,
	}

	return snowflake
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置数据位
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) setBits() {
	s.nodeMax = -1 ^ (-1 << s.nodeBits)
	s.nodeMask = s.nodeMax << s.seqBits
	s.seqMask = -1 ^ (-1 << s.seqBits)

	s.timeShift = s.nodeBits + s.seqBits
	s.nodeShift = s.seqBits
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置开始时间戳
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) SetStartTimestamp(startTimestamp int64) {
	if startTimestamp <= 0 {
		startTimestamp = 1551377400197
	}

	s.startTimestamp = startTimestamp
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置节点数据位数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) SetNodeBits(nodeBits uint8) {
	if nodeBits == 0 {
		nodeBits = 10
	}

	s.nodeBits = nodeBits
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取节点数据位数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) GetNodeBits() uint8 {
	return s.nodeBits
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置序列号数据位数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) SetSeqBits(seqBits uint8) {
	if seqBits == 0 {
		seqBits = 12
	}

	s.seqBits = seqBits
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取序列号数据位数
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) GetSeqBits() uint8 {
	return s.seqBits
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取节点对象
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Snowflake) GetNode(args ...int64) *SnowflakeNode {
	s.setBits()

	var workerNodeId int64
	if len(args) > 0 {
		workerNodeId = args[0]
	}

	if workerNodeId < 0 || workerNodeId > s.nodeMax {
		workerNodeId = 0
	}

	snowflakeNode := &SnowflakeNode{
		snowflake:    s,
		workerNodeId: workerNodeId,
		seq:          0,
		timestamp:    0,
	}

	return snowflakeNode
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Id数值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (n *SnowflakeNode) GetId() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	//纳秒转换成毫秒
	now := time.Now().UnixNano() / 1000000

	if n.timestamp == now {
		//同一毫秒时间周期内，序号递增
		n.seq = (n.seq + 1) & n.snowflake.seqMask

		//序号已超过最大值归零，则等待下一个时间周期
		if n.seq == 0 {
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		//新时间周期序号初始化为0
		n.seq = 0
	}

	n.timestamp = now

	return int64((now-n.snowflake.startTimestamp)<<n.snowflake.timeShift |
		(n.workerNodeId << n.snowflake.nodeShift) |
		(n.seq),
	)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Id数值的时间戳值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (n *SnowflakeNode) GetTimestamp(id int64) int64 {
	return (id >> n.snowflake.timeShift) + n.snowflake.startTimestamp
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Id数值的节点值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (n *SnowflakeNode) GetNode(id int64) int64 {
	return id & n.snowflake.nodeMask >> n.snowflake.nodeShift
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Id数值的序号值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (n *SnowflakeNode) GetSeq(id int64) int64 {
	return id & n.snowflake.seqMask
}
