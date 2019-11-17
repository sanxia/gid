package gid

import (
	"sync"
)

/* ================================================================================
 * gid
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

type (
	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * id data source interface
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	IIdSource interface {
		GetNextId() uint64
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * id store interface
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	IIdStore interface {
		LoadPrevId() uint64
		SavePrevId(id uint64)
	}
)

type (
	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * gid data struct
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	gid struct {
		count          int
		threshold      int
		prevId         uint64
		nextId         uint64
		totalIndex     uint64
		currentIndex   int
		frontDataChan  chan NextData
		backDataChan   chan NextData
		frontCacheChan chan bool
		backCacheChan  chan bool
		dataSource     IIdSource
		dataStore      IIdStore
		idMux          sync.Mutex
		dataMux        sync.Mutex
		isFront        bool
	}

	NextData struct {
		id     uint64
		isLast bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instance gid data struct
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewGid(count, threshold int) *gid {
	generationId := &gid{
		frontDataChan:  make(chan NextData, count),
		backDataChan:   make(chan NextData, count),
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
func (s *gid) init() {
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
 * load id data
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *gid) Start() {
	if s.dataStore != nil {
		s.prevId = s.dataStore.LoadPrevId()
	}

	s.frontCacheChan <- true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get next id
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *gid) GetNextId() uint64 {
	s.idMux.Lock()
	defer s.idMux.Unlock()

	nextId := s.getNextId()
	for nextId <= s.prevId {
		nextId = s.getNextId()
	}

	s.nextId = nextId

	return s.nextId
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set IIdSource interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *gid) SetSource(dataSource IIdSource) {
	s.dataSource = dataSource
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set IIdStore interface
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *gid) SetStore(dataStore IIdStore) {
	s.dataStore = dataStore
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * getNextId
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *gid) getNextId() uint64 {
	var nextData NextData

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

	if nextData.isLast {
		s.currentIndex = 0
		s.isFront = !s.isFront

		if s.dataStore != nil {
			s.prevId = s.dataStore.LoadPrevId()
			if nextData.id > s.prevId {
				s.dataStore.SavePrevId(nextData.id)
			}
		}
	}

	return nextData.id
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * loadData
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *gid) loadData(dataChan chan NextData) {
	s.dataMux.Lock()
	defer s.dataMux.Unlock()

	for i := 0; i < s.count; i++ {
		s.totalIndex = s.totalIndex + 1

		isLast := false
		if i == s.count-1 {
			isLast = true
		}

		nextId := s.totalIndex
		if s.dataSource != nil {
			nextId = s.dataSource.GetNextId()
		}

		nextData := NextData{
			id:     nextId,
			isLast: isLast,
		}

		dataChan <- nextData
	}
}
