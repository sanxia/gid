package gid

import (
	"sync"
)

/* ================================================================================
 * id generator buffer
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

type (
	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * id data source interface
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	IIdSource interface {
		GetNextId() int64
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * id store interface
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	IIdStore interface {
		LoadPrevId() int64
		SavePrevId(id int64)
	}
)

type (
	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * IdGenerator data struct
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	IdGenerator struct {
		count          int
		threshold      int
		prevId         int64
		nextId         int64
		totalIndex     int64
		currentIndex   int
		frontDataChan  chan IdData
		backDataChan   chan IdData
		frontCacheChan chan bool
		backCacheChan  chan bool
		dataSource     IIdSource
		dataStore      IIdStore
		idMux          sync.Mutex
		dataMux        sync.Mutex
		isFront        bool
	}

	IdData struct {
		id     int64
		isHead bool
		isTail bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instance IdGenerator data struct
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewIdGenerator(args ...int) *IdGenerator {
	count := 1000
	if len(args) > 0 {
		count = args[0]
	}

	threshold := int(float64(count) * (float64(10) / float64(100)))
	if threshold <= 0 {
		threshold = 1
	}

	generationId := &IdGenerator{
		frontDataChan:  make(chan IdData, count),
		backDataChan:   make(chan IdData, count),
		frontCacheChan: make(chan bool),
		backCacheChan:  make(chan bool),
		count:          count,
		threshold:      threshold,
		isFront:        true,
	}

	generationId.init()

	return generationId
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * cache chan listen
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) init() {
	go func() {
		for {
			select {
			case <-s.frontCacheChan:
				s.loadData(s.frontDataChan)
			case <-s.backCacheChan:
				s.loadData(s.backDataChan)
			}
		}
	}()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * launch data
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) Launch() {
	if s.dataStore != nil {
		s.prevId = s.dataStore.LoadPrevId()
	}

	s.frontCacheChan <- true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get next id
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) GetNextId() int64 {
	s.idMux.Lock()
	defer s.idMux.Unlock()

	nextData := s.getNextData()
	for nextData.id <= s.prevId {
		nextData = s.getNextData()
	}

	s.nextId = nextData.id

	//switch buffer
	if nextData.isTail {
		s.currentIndex = 0
		s.isFront = !s.isFront
	}

	return s.nextId
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set count
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) SetCount(count int) {
	if count <= 0 {
		count = 1000
	}

	s.count = count

	s.threshold = int(float64(count) * (float64(10) / float64(100)))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set IIdSource interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) SetSource(dataSource IIdSource) {
	s.dataSource = dataSource
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set IIdStore interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) SetStore(dataStore IIdStore) {
	s.dataStore = dataStore
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * getNextId
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) getNextData() IdData {
	var nextData IdData

	if s.isFront {
		select {
		case nextData = <-s.frontDataChan:
			s.currentIndex = s.currentIndex + 1

			if s.count-s.currentIndex == s.threshold {
				s.backCacheChan <- true
			}
		}
	} else {
		select {
		case nextData = <-s.backDataChan:
			s.currentIndex = s.currentIndex + 1

			if s.count-s.currentIndex == s.threshold {
				s.frontCacheChan <- true
			}
		}
	}

	return nextData
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * loadData
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *IdGenerator) loadData(dataChan chan IdData) {
	s.dataMux.Lock()
	defer s.dataMux.Unlock()

	for i := 0; i < s.count; i++ {
		s.totalIndex = s.totalIndex + 1

		isHead := false
		isTail := false
		if i == 0 {
			isHead = true
		} else if i == s.count-1 {
			isTail = true
		}

		//next data
		nextId := s.totalIndex
		if s.dataSource != nil {
			nextId = s.dataSource.GetNextId()
		}

		nextData := IdData{
			id:     nextId,
			isHead: isHead,
			isTail: isTail,
		}

		//report store
		if isTail {
			if s.dataStore != nil {
				s.prevId = s.dataStore.LoadPrevId()
				if nextData.id > s.prevId {
					s.dataStore.SavePrevId(nextData.id)
				}
			}
		}

		dataChan <- nextData
	}
}
