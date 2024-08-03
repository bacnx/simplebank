package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bacnx/simplebank/api"
	db "github.com/bacnx/simplebank/db/sqlc"
	_ "github.com/bacnx/simplebank/doc/statik"
	"github.com/bacnx/simplebank/gapi"
	"github.com/bacnx/simplebank/mail"
	"github.com/bacnx/simplebank/pb"
	"github.com/bacnx/simplebank/util"
	"github.com/bacnx/simplebank/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.GetConfig(".")
	if err != nil {
		log.Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisDistributor(redisOpt)

	go runTaskProcessor(config, redisOpt, store)
	go runRrpcGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)
}

func runDBMigration(migrationUrl, dbSource string) {
	m, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Err(err).Msg("cannot create migrate instance")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Err(err).Msg("cannot run migrate up")
	}
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisConnOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.FromEmailName, config.FromEmailAddress, config.FromEmailPassword)
	taskProcessor := worker.NewRedisProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Err(err)
	}

	server.Start(config.HTTPServerAddress)
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Err(err).Msg("cannot create server")
	}

	serverOption := grpc.ChainUnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(serverOption)
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Err(err).Msg("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Err(err).Msg("cannot start gRPC server")
	}
}

func runRrpcGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Err(err).Msg("cannot create gapi server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSimplebankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	handler := gapi.HttpLogger(grpcMux)
	mux.Handle("/", handler)

	statikFS, err := fs.New()
	if err != nil {
		log.Err(err).Msg("cannot crate file system")
	}

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Err(err).Msg("cannot create listener")
	}

	log.Printf("start http server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Err(err).Msg("cannot start http server")
	}
}
