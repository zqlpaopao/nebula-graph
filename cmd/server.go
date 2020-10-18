package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	v1 "test/grpc/api/service/v1"
	"test/grpc/server"
)

type Config struct {
	GRPCPort             string
	DatasStordDBHost     string
	DatasStordDBUser     string
	DatasStordDBPassword string
	DatasStordDBSchema   string
}

func RunServer() error {
	ctx := context.Background()
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatasStordDBHost, "db-host", "", "db-host")
	flag.StringVar(&cfg.DatasStordDBUser, "db-user", "", "db-user")
	flag.StringVar(&cfg.DatasStordDBPassword, "db-password", "", "db-password")
	flag.StringVar(&cfg.DatasStordDBSchema, "db-schema", "", "db-schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server:%s", cfg.GRPCPort)
	}
	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.DatasStordDBUser, cfg.DatasStordDBPassword, cfg.DatasStordDBHost, cfg.DatasStordDBSchema, param)

	db, err := sql.Open("mysql", dsn)
	if nil != err {
		return fmt.Errorf("连接数据库失败:%v", err)
	}
	defer db.Close()
	v1API := v1.NewToDoServiceServer(db)
	return server.RunServer(ctx, v1API, cfg.GRPCPort)
}
