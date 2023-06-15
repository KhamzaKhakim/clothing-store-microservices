package main

import (
	"clothing-store-brands/internal/data"
	"clothing-store-brands/internal/server"
	"clothing-store-brands/pkg/config"
	"clothing-store-brands/pkg/jsonlog"
	pb "clothing-store-brands/pkg/pb"
	"context"
	"database/sql"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"os"
	"time"
	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	startGrpc()
	connGrpc, err := grpc.Dial("localhost:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer connGrpc.Close()

	client := pb.NewBrandsServiceClient(connGrpc)

	//-------------------------------
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	chSingle, err := conn.Channel()
	failOnError(err, "Failed to open a channel to send message to single client")
	defer chSingle.Close()

	q, err := chSingle.QueueDeclare(
		"brand_queue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = chSingle.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := chSingle.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {

		ctx := context.Background()

		for d := range msgs {

			req := &pb.ShowBrandRequest{}

			proto.Unmarshal(d.Body, req)

			res, err := client.ShowBrand(ctx, req)

			responseByte, err := proto.Marshal(res)
			if err != nil {
				log.Fatalf("could not marshal: %v", err)
			}
			//
			err = chSingle.PublishWithContext(ctx,
				"",      // exchange
				"brand", // routing key
				false,   // mandatory
				false,   // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        responseByte,
				})

			failOnError(err, "Failed to publish a message to user")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.C.Db.Dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.C.Db.MaxOpenConns)
	db.SetMaxIdleConns(config.C.Db.MaxIdleConns)
	duration, err := time.ParseDuration(config.C.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func startGrpc() {
	config.ReadConfig()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	logger.PrintInfo("database connection pool established", nil)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.C.Port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterBrandsServiceServer(s, &server.Server{
		Logger: logger,
		Models: data.NewModels(db),
	})

	log.Printf("Server listening at %d", config.C.Port)

	go startServer(s, lis)
}

func startServer(s *grpc.Server, lis net.Listener) {
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
