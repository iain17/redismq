package redismq

import (
	"github.com/adeven/goenv"
	"github.com/adeven/redis"
)

type Queue struct {
	redisClient *redis.Client
	Name        string
}

func NewQueue(goenv *goenv.Goenv, name string) *Queue {
	q := &Queue{Name: name}
	host, port, db := goenv.GetRedis()
	q.redisClient = redis.NewTCPClient(host+":"+port, "", int64(db))
	q.redisClient.SAdd(MasterQueueKey(), name)
	return q
}

func MasterQueueKey() string {
	return "redismq::queues"
}

func (self *Queue) InputName() string {
	return "redismq::" + self.Name
}

func (self *Queue) WorkingName(consumer string) string {
	return "redismq::" + self.Name + "::working::" + consumer
}

func (self *Queue) FailedName() string {
	return "redismq::" + self.Name + "::failed"
}

func (self *Queue) InputCounterName() string {
	return self.InputName() + "::counter"
}

func (self *Queue) WorkingCounterName(consumer string) string {
	return self.WorkingName(consumer) + "::counter"
}

func (self *Queue) FailedCounterName() string {
	return self.FailedName() + "::counter"
}

func (self *Queue) AckCounterName(consumer string) string {
	return self.InputName() + "::ack::" + consumer + "::counter"
}

func (self *Queue) HasUnacked(consumer string) bool {
	if self.GetUnackedLength(consumer) != 0 {
		return true
	}
	return false
}

func (self *Queue) parseRedisAnswer(answer *redis.StringReq, consumer string) (*Package, error) {
	if answer.Err() != nil {
		return nil, answer.Err()
	}
	p, err := UnmarshalPackage(answer.Val(), self, consumer)
	if err != nil {
		return nil, err
	}
	return p, nil
}