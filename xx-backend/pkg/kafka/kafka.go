package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
	config *Config
}

type Config struct {
	Brokers         []string
	TopicUserEvents string
	TopicSystemLogs string
}

type Message struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func NewKafkaClient(config *Config) (*KafkaClient, error) {
	// 创建生产者
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	// 创建消费者
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Brokers,
		Topic:    config.TopicUserEvents, // 默认监听用户事件主题
		GroupID:  "xx-backend-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaClient{
		writer: writer,
		reader: reader,
		config: config,
	}, nil
}

// SendUserEvent 发送用户事件
func (kc *KafkaClient) SendUserEvent(eventType string, data map[string]interface{}) error {
	message := Message{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = kc.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: kc.config.TopicUserEvents,
		Key:   []byte(eventType),
		Value: jsonData,
	})

	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Sent user event: %s", eventType)
	return nil
}

// SendSystemLog 发送系统日志
func (kc *KafkaClient) SendSystemLog(logType string, data map[string]interface{}) error {
	message := Message{
		Type:      logType,
		Timestamp: time.Now(),
		Data:      data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = kc.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: kc.config.TopicSystemLogs,
		Key:   []byte(logType),
		Value: jsonData,
	})

	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Sent system log: %s", logType)
	return nil
}

// ConsumeMessages 消费消息
func (kc *KafkaClient) ConsumeMessages(ctx context.Context, handler func(Message) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var message Message
			if err := json.Unmarshal(m.Value, &message); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}

// ConsumeUserEvents 消费用户事件
func (kc *KafkaClient) ConsumeUserEvents(ctx context.Context, handler func(Message) error) error {
	// 切换到用户事件主题
	kc.reader.SetOffset(kafka.LastOffset)

	return kc.ConsumeMessages(ctx, handler)
}

// ConsumeSystemLogs 消费系统日志
func (kc *KafkaClient) ConsumeSystemLogs(ctx context.Context, handler func(Message) error) error {
	// 创建系统日志消费者
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kc.config.Brokers,
		Topic:    kc.config.TopicSystemLogs,
		GroupID:  "xx-backend-system-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading system log: %v", err)
				continue
			}

			var message Message
			if err := json.Unmarshal(m.Value, &message); err != nil {
				log.Printf("Error unmarshaling system log: %v", err)
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Error handling system log: %v", err)
			}
		}
	}
}

// Close 关闭Kafka连接
func (kc *KafkaClient) Close() error {
	if err := kc.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	if err := kc.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	return nil
}

// CreateTopics 创建主题
func (kc *KafkaClient) CreateTopics() error {
	conn, err := kafka.Dial("tcp", kc.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	topics := []string{kc.config.TopicUserEvents, kc.config.TopicSystemLogs}

	for _, topic := range topics {
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		}

		err = conn.CreateTopics(topicConfigs...)
		if err != nil {
			// 如果主题已存在，忽略错误
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("Topic %s already exists", topic)
				continue
			}
			return fmt.Errorf("failed to create topic %s: %w", topic, err)
		}
		log.Printf("Created topic: %s", topic)
	}

	return nil
}
