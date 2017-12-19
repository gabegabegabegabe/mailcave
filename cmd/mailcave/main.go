package main

import (
	"context"
	"flag"
	"time"

	"github.com/tambchop/mailcave/logging"
	"github.com/tambchop/mailcave/mail"
)

func main() {

	dbAddr := flag.String("dbAddr", "mongodb://localhost:27017", "the address of the mailcave database")
	dbName := flag.String("dbName", "mailcave", "the name of the mailcave archive database")
	ipAddr := flag.String("ipAddr", ":8080", "the ip address that mailcave should listen on")
	logToStdOut := flag.Bool("logToStdOut", true, "should we also log to stdout in addition to file?")
	flag.Parse()

	logger := logging.NewLogger("mailcave", "./logs", 25, 10, 30, *logToStdOut)
	defer logger.Close()

	mongo := mail.NewMongoArchive(*dbAddr, *dbName, logger)
	archivist := mail.NewArchivist(mongo, logger)

	time.Sleep(1 * time.Second)

	ctx := context.Background()
	err := archivist.Start(ctx, *ipAddr)
	if err != nil {
		logger.Printf("failed to start Archiver with error '%s'", err)
	}
}
