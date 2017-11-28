package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/lightstaff/go-dddcqrses/command"
	"github.com/lightstaff/go-dddcqrses/common"
	"github.com/lightstaff/go-dddcqrses/messages"
	"github.com/lightstaff/go-dddcqrses/query"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	inMemoryDB := common.NewInMemoryDB(sugar)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	delivery := make(chan *common.InMemoryMessage)
	errc := make(chan error)

	go func() {
		persistence := common.NewFakePersistence(inMemoryDB, sugar)
		producer := common.NewFakeMessagingProducer(delivery, sugar)
		commandActor := command.NewTodoActor(persistence, producer, sugar)
		commandActor.Act(&command.TodoRegistry{
			Message:   "test message",
			Completed: false,
		})
		sugar.Info("command actor action completed")
	}()

	go func() {
		defer close(delivery)

		persistenceQuery := common.NewFakePersistenceQuery(inMemoryDB, sugar)
		consumer := common.NewFakeMessagingConsumer(delivery, sugar)
		queryDB := query.NewFakeQueryDB(sugar)
		queryActor := query.NewTodoActor(persistenceQuery, queryDB, sugar)

		go func() {
			if err := consumer.Consume(ctx, delivery); err != nil {
				errc <- err
			}
		}()

		for {
			select {
			case <-ctx.Done():
				sugar.Info("call context cancel")
				return
			case m := <-delivery:
				msg, err := messages.MessageConveter(m.Header, m.Data)
				if err != nil {
					errc <- err
				}
				sugar.Infow("receive message", "message", msg)
				switch msg := msg.(type) {
				case *messages.TodoEventOccurred:
					if err := queryActor.Act(msg); err != nil {
						errc <- err
					}
					sugar.Info("query actor action completed")
					sugar.Infow("now todo", "todo", queryDB.FindByID(msg.AggregateID))
				}
			}
		}
	}()

	go func() {
		if err := <-errc; err != nil {
			sugar.Errorw("error happend", "error", err)
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
