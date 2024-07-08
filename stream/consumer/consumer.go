package consumer

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func CheckErrReceive(err error) {
	if err != nil {
		fmt.Printf("%s ", err)
		os.Exit(1)
	}
}
func main() {

	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"))
	CheckErrReceive(err)

	streamName := "hello-go-stream"
	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	CheckErrReceive(err)

	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		fmt.Printf("Поток: %s - Сообщение: %s\n", consumerContext.Consumer.GetStreamName(),
			message.Data)
	}

	consumer, err := env.NewConsumer(streamName, messagesHandler,
		stream.NewConsumerOptions().SetOffset(stream.OffsetSpecification{}.First()))
	CheckErrReceive(err)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(" [x] Ожидание сообщения. enter для закрытия получателя")
	_, _ = reader.ReadString('\n')
	err = consumer.Close()
	CheckErrReceive(err)

}
