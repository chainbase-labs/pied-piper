package kafka

import (
	"time"

	"github.com/ddl-hust/pied-piper/common/log"

	"github.com/ddl-hust/pied-piper/common/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"go.uber.org/zap"
)

func NewClient(topic string) (*kgo.Client, error) {
	cfg := config.GetConf()
	opts := make([]kgo.Opt, 0)

	opts = append(opts, kgo.SeedBrokers(cfg.KafkaHosts))
	if cfg.KafkaUser != "" && cfg.KafkaPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{
			User: cfg.KafkaUser,
			Pass: cfg.KafkaPass,
		}.AsSha256Mechanism()))
	}
	opts = append(opts, kgo.ConsumerGroup(cfg.ConsumeGroup))
	opts = append(opts, kgo.ConsumeTopics(topic))
	opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AfterMilli(time.Now().UnixMilli())))

	kafkaClient, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatal("failed to create kafka client", zap.Error(err))
		return nil, err
	}
	log.Info("create kafka client success")
	return kafkaClient, nil
}
