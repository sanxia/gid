# gid
double buffer generation id

===== Example =====

package main

import (
    "log"

    "time"
)

import (

    "github.com/sanxia/gid"

)

type IdData struct {

}

// custom extend get next id 
func (s *IdData) GetNextId() uint64 {

    time.Sleep(3 * time.Second)

    return 1001
}

// custom extend load prev id 
// example: load for redis
func (s *IdData) LoadPrevId() uint64 {

    return 1000

}

// custom extend save prev id 
// example: save for redis
func (s *IdData) SavePrevId(prevId uint64) {

}

func main() {

    idData := new(IdData)

    idGeneration := gid.NewGid(10, 5)

    //idGeneration.SetSource(idData)

    idGeneration.SetStore(idData)

    idGeneration.Start()

    for i := 0; i < 2; i++ {

        go func(currentIndex int) {

            for {

                log.Printf("go routine %02d %v NextId = %d", currentIndex, time.Now(), idGeneration.GetNextId())

            }

        }(i)

    }

    for {

        log.Printf("go routine 03 %v NextId = %d", time.Now(), idGeneration.GetNextId())

    }

}

