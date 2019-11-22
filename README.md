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

type IdProvider struct {

    snowflakeNode *gid.SnowflakeNode

}

func NewIdProvider(snowflakeNode *gid.SnowflakeNode) *IdProvider {

    return &IdProvider{

        snowflakeNode: snowflakeNode,

    }

}

func (s *IdProvider) GetNextId() int64 {

    return s.snowflakeNode.GetId()

}

func (s *IdProvider) LoadPrevId() int64 {

    return 0

}

func (s *IdProvider) SavePrevId(prevId int64) {

    log.Printf("SavePrevId: %d", prevId)

}


func main() {

    snowflake := gid.NewSnowflake()

    snowflake.SetStartTimestamp(1571145696789)

    snowflake.SetNodeBits(10)

    snowflake.SetSeqBits(12)

    idProvider := NewIdProvider(snowflake.GetNode(1))

    idGenerator := gid.NewIdGenerator()

    idGenerator.SetCount(5000)

    idGenerator.SetSource(idProvider)

    idGenerator.SetStore(idProvider)

    idGenerator.Launch()

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





