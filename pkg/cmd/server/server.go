package cmd

import (
	"context"
	"flag"
	"fmt"
    "database/sql"

	"github.com/dansusman/todoservice/pkg/protocol/grpc"

    v1 "github.com/dansusman/todoservice/pkg/service/v1"
)

type Config struct {
    GRPCPort string
    DBHost string
    DBUser string
    DBPassword string
    DBSchema string
}

func RunServer() error {
    ctx := context.Background()

    var config Config

    flag.StringVar(&config.GRPCPort, "grpc-port", "", "gRPC port to bind")
    flag.StringVar(&config.DBHost, "db-host", "", "Database host")
    flag.StringVar(&config.DBUser, "db-user", "", "Database user")
    flag.StringVar(&config.DBPassword, "db-password", "", "Database password")
    flag.StringVar(&config.DBSchema, "db-schema", "", "Database schema")
    flag.Parse()

    if len(config.GRPCPort) == 0 {
        return fmt.Errorf("invalid TCP port for gRPC server: '%s'", config.GRPCPort)
    }

    // add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBSchema,
		param)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()


    v1Api := v1.NewTodoServiceServer(db)

    return grpc.RunServer(ctx, v1Api, config.GRPCPort)
}
