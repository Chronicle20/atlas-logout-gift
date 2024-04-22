package main

import (
	"atlas-lgs/database"
	"atlas-lgs/logger"
	"atlas-lgs/rest"
	"atlas-lgs/tracing"
	"atlas-lgs/world"
	"atlas-lgs/world/character/choice"
	"atlas-lgs/world/gift"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-lgs"

type Server struct {
	baseUrl string
	prefix  string
}

func (s Server) GetBaseURL() string {
	return s.baseUrl
}

func (s Server) GetPrefix() string {
	return s.prefix
}

func GetServer() Server {
	return Server{
		baseUrl: "",
		prefix:  "/api/lgs/",
	}
}

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	db := database.Connect(l, database.SetMigrations(gift.Migration, choice.Migration))

	rest.CreateService(l, db, ctx, wg, GetServer().GetPrefix(), gift.InitResource(GetServer()), choice.InitResource(GetServer()))

	initializeGifts(l, db)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}

func initializeGifts(l logrus.FieldLogger, db *gorm.DB) {
	s := gift.GetAll(l, db)
	if len(s) > 0 {
		return
	}

	filePath, ok := os.LookupEnv("SEED_JSON_FILE_PATH")
	if !ok {
		l.Fatalf("Environment variable SEED_JSON_FILE_PATH is not set.")
	}

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		l.Fatal("Error reading JSON file:", err)
	}

	// Define a slice to store the objects
	var objects []world.JSONModel
	var gifts []gift.Model

	// Unmarshal JSON into the slice
	err = json.Unmarshal(jsonData, &objects)
	if err != nil {
		l.Fatal("Error unmarshalling JSON:", err)
	}

	for _, jdo := range objects {
		for _, gdo := range jdo.Gifts {
			md := gift.NewGiftBuilder(0).
				SetWorldId(jdo.Id).
				SetItemId(gdo.ItemId).
				SetSerialNumber(gdo.SerialNumber).
				SetWeight(gdo.Weight).
				Build()
			gifts = append(gifts, md)
		}
	}

	err = gift.BulkCreateGifts(db, gifts)
	if err != nil {
		l.Fatalf(err.Error())
	}
}
