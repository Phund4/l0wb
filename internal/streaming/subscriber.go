package streaming

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/Phund4/l0wb/internal/db"

	stan "github.com/nats-io/stan.go"
)

type Subscriber struct {
	sub      stan.Subscription
	dbObject *db.DB
	sc       *stan.Conn
	name     string
}

// Создание нового Subscriber
func NewSubscriber(db *db.DB, conn *stan.Conn) *Subscriber {
	return &Subscriber{
		name:     "Subscriber",
		dbObject: db,
		sc:       conn,
	}
}

// Подписка Subscriber
func (s *Subscriber) Subscribe() {
	var err error

	ackWait, err := strconv.Atoi(os.Getenv("NATS_ACK_WAIT_SECONDS"))
	if err != nil {
		log.Printf("%s: received a message!\n", s.name)
		return
	}

	s.sub, err = (*s.sc).Subscribe(
		os.Getenv("NATS_SUBJECT"),
		func(m *stan.Msg) {
			log.Printf("%s: received a message!\n", s.name)
			if s.messageHandler(m.Data) {
				err := m.Ack()
				if err != nil {
					log.Printf("%s ack() err: %s", s.name, err)
				}
			}
		},
		stan.AckWait(time.Duration(ackWait)*time.Second),
		stan.DurableName(os.Getenv("NATS_DURABLE_NAME")),
		stan.SetManualAckMode(),
		stan.MaxInflight(10))
		
	if err != nil {
		log.Printf("%s: error: %v\n", s.name, err)
	}
	log.Printf("%s: subscribed to subject %s\n", s.name, os.Getenv("NATS_SUBJECT"))
}

// Парсинг полученных данных и отправка в DB и Cache
func (s *Subscriber) messageHandler(data []byte) bool {
	recievedOrder := db.Order{}
	err := json.Unmarshal(data, &recievedOrder)
	if err != nil {
		log.Printf("%s: messageHandler() error, %v\n", s.name, err)
		return true
	}
	log.Printf("%s: unmarshal Order to struct: %v\n", s.name, recievedOrder)

	_, err = s.dbObject.AddOrder(recievedOrder)
	if err != nil {
		log.Printf("%s: unable to add order to db: %v\n", s.name, err)
		return false
	}

	err = s.dbObject.UploadCache();
	if err != nil {
		log.Printf("%s: unable to add order to cache: %v\n", s.name, err)
		return false
	}

	return true
}

// Отписка
func (s *Subscriber) Unsubscribe() {
	if s.sub != nil {
		s.sub.Unsubscribe()
	}
}
