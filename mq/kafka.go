package mq

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"time"
)

var (
	producer sarama.SyncProducer
	consumer sarama.Consumer
)

func InitKafka(adders []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	//config.producer.Partitioner = sarama.NewRandomPartitioner	// 随机分区
	//config.producer.Partitioner = sarama.NewHashPartitioner	// hash分区
	config.Producer.Partitioner = sarama.NewManualPartitioner // 可以用sarama.ProducerMessage的Partition，指定分区
	//config.producer.Partitioner = sarama.NewRoundRobinPartitioner // 循环分区

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest //初始从最新的offset开始

	c, err := sarama.NewClient(adders, config)
	if err != nil {
		panic(err)
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	producer = p

	cs, err := sarama.NewConsumerFromClient(c)
	if err != nil {
		panic(err)
	}
	consumer = cs
}

func Publish(topic string, msg interface{}) error {
	data, _ := json.Marshal(msg)

	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(data),
	})

	return err
}

func Subscribe(topic string, f func(key string, val []byte)) {
	pl, err := consumer.Partitions(topic)

	if err != nil {
		panic(err)
	}
	for _, p := range pl {
		cp, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		go func(_cp sarama.PartitionConsumer) {
			for msg := range _cp.Messages() {
				f(string(msg.Key), msg.Value)
			}
		}(cp)
	}
}
