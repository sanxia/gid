# gid
double buffer generation id

# Example 

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

# Run Result

go routine 01 2019-11-17 23:00:12.23975 +0800 CST m=+3.086219504 NextId = 102078

go routine 00 2019-11-17 23:00:12.238863 +0800 CST m=+3.085331738 NextId = 102050

go routine 00 2019-11-17 23:00:12.239776 +0800 CST m=+3.086244983 NextId = 102080

go routine 00 2019-11-17 23:00:12.239792 +0800 CST m=+3.086261411 NextId = 102081

go routine 03 2019-11-17 23:00:12.239021 +0800 CST m=+3.085490485 NextId = 102059

go routine 00 2019-11-17 23:00:12.239809 +0800 CST m=+3.086277883 NextId = 102082

go routine 00 2019-11-17 23:00:12.239823 +0800 CST m=+3.086292095 NextId = 102084

go routine 00 2019-11-17 23:00:12.239836 +0800 CST m=+3.086304757 NextId = 102085

go routine 00 2019-11-17 23:00:12.239981 +0800 CST m=+3.086450238 NextId = 102086

go routine 00 2019-11-17 23:00:12.240003 +0800 CST m=+3.086471830 NextId = 102087

go routine 00 2019-11-17 23:00:12.240022 +0800 CST m=+3.086490642 NextId = 102088





