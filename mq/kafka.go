package mq

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"time"
)

type Kafka struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

var kafkaList = make(map[string]*Kafka)

func ConnectKafka(name string, adders []string) error {
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
		return err
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return err
	}

	cs, err := sarama.NewConsumerFromClient(c)
	if err != nil {
		return err
	}

	kafkaList[name] = &Kafka{
		producer: p,
		consumer: cs,
	}

	return nil
}

func Publish(name, topic string, msg interface{}) error {
	kf := kafkaList[name]

	data, _ := json.Marshal(msg)
	_, _, err := kf.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(data),
	})

	return err
}

func Subscribe(name, topic string, f func(key string, val []byte)) error {
	kf := kafkaList[name]
	consumer := kf.consumer

	pl, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for _, p := range pl {
		cp, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(_cp sarama.PartitionConsumer) {
			for msg := range _cp.Messages() {
				f(string(msg.Key), msg.Value)
			}
		}(cp)
	}

	return nil
}
