package streaming

import (
	"encoding/json"
	"log"
	"os"
	"math/rand"
	"time"

	"github.com/Phund4/l0wb/internal/db"

	stan "github.com/nats-io/stan.go"
)

type Publisher struct {
	sc   *stan.Conn
	name string
}

func NewPublisher(conn *stan.Conn) *Publisher {
	return &Publisher{
		name: "Publisher",
		sc:   conn,
	}
}

// Тестовый скрипт публикации данных Order
func (p *Publisher) Publish() {
	order := createOrder();
	orderData, err := json.Marshal(order)
	if err != nil {
		log.Printf("%s: json.Marshal error: %v\n", p.name, err)
	}

	// An asynchronous publish API
	ackHandler := func(ackedNuid string, err error) {
		if err != nil {
			log.Printf("%s: error publishing msg id %s: %v\n", p.name, ackedNuid, err.Error())
		} else {
			log.Printf("%s: received ack for msg id: %s\n", p.name, ackedNuid)
		}
	}

	// публикация данных:
	log.Printf("%s: publishing data ...\n", p.name)
	nuid, err := (*p.sc).PublishAsync(os.Getenv("NATS_SUBJECT"), orderData, ackHandler)
	if err != nil {
		log.Printf("%s: error publishing msg %s: %v\n", p.name, nuid, err.Error())
	}
}

func randomStr(str string, strLen int) string {
	var res string
	for i := 0; i < strLen; i++ {
		ch := string(str[rand.Intn(strLen)])
		res += ch
	}
	return res
}

func createOrder() db.Order {
	letters := "abcdefghijklmnopqrstuvwxyz"
	numbers := "1234567890"
	lettersAndNumbers := letters + numbers;

	name := randomStr(letters, 6)
	delivery := db.Delivery{
		Name:    name + " " + name + "ov",
		Phone:   "+97" + randomStr(numbers, 8),
		Zip:     randomStr(numbers, 3),
		City:    randomStr(letters, 6),
		Address: randomStr(letters, 5) + randomStr(numbers, 2),
		Region:  randomStr(letters, 6),
		Email:   randomStr(lettersAndNumbers, 7) + "@gmail.com",
	}

	payment := db.Payment{
		Transaction:  randomStr(lettersAndNumbers, 6),
		RequestID:    "",
		Currency:     randomStr(letters, 6),
		Provider:     randomStr(letters, 6),
		Amount:       rand.Intn(10000),
		PaymentDt:    int64(rand.Intn(9000)),
		Bank:         randomStr(letters, 6),
		DeliveryCost: rand.Intn(9000),
		GoodsTotal:   rand.Intn(900),
		CustomFee:    rand.Intn(10000),
	}

	item := db.Item{
		ChrtID:      rand.Intn(100000),
		TrackNumber: randomStr(letters, 5),
		Price:       rand.Intn(1000),
		RID:         randomStr(lettersAndNumbers, 6),
		Name:        randomStr(letters, 6),
		Sale:        rand.Intn(100),
		Size:        randomStr(numbers, 2),
		TotalPrice:  rand.Intn(100000),
		NmID:        rand.Intn(90000),
		Brand:       randomStr(letters, 6),
		Status:      rand.Intn(100),
	}
	
	timeOrder := time.Now();
	order := db.Order{
		TrackNumber:       randomStr(letters, 6),
		Entry:             randomStr(letters, 6),
		Locale:            randomStr(letters, 6),
		InternalSignature: "",
		CustomerID:        randomStr(letters, 6),
		DeliveryService:   randomStr(letters, 6),
		ShardKey:          randomStr(numbers, 4),
		SmID:              rand.Intn(100),
		DateCreated:       timeOrder,
		OofShard:          randomStr(numbers, 4),
	}

	order.OrderUID = randomStr(numbers, 7);
	order.Delivery = delivery
	order.Payment = payment
	order.Items = append(order.Items, item)

	return order
}