package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	orderFields = `OrderUID, TrackNumber, Entry, Locale, InternalSignature, CustomerID, DeliveryService, Shardkey, SmID, 
	DateCreated, OofShard`
	deliveryFields = `Name, Phone, Zip, City, Address, Region, Email`
	paymentFields  = `Transaction, RequestID, Currency, Provider, Amount, PaymentDt, Bank, DeliveryCost, GoodsTotal, CustomFee`
	itemFields     = `ChrtID, TrackNumber, Price, RID, Name, Sale, Size, TotalPrice, NmID, Brand, Status`
)

type DB struct {
	sqlxDB *sqlx.DB
	cache  *Cache
	name   string
}

func NewDB(csh *Cache) *DB {
	db := DB{cache: csh}
	db.Init()
	return &db
}

func (db *DB) GetOrderFromCache(OrderUid string) (Order, error) {
	log.Printf("Service is trying to find the order %s...\n", OrderUid)
	return db.cache.GetOrder(OrderUid)
}

// Получение Order из БД по id
func (db *DB) GetOrderFromDB(oid int64) (Order, error) {
	var o Order
	tx, err := db.sqlxDB.Beginx()
	if err != nil {
		log.Printf("%s: Error in init sqlxDB: %v", db.name, err)
		return o, err
	}

	orderQuery := fmt.Sprintf(`SELECT %s FROM orders WHERE id = $1`, orderFields)
	// Сбор данных об Order
	err = tx.QueryRow(orderQuery, oid).Scan(&o.OrderUID, &o.TrackNumber, &o.Entry,
		&o.Locale, &o.InternalSignature, &o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard)
	if err != nil {
		log.Printf("%v: unable to get order from database: %v\n", db.name, err)
		return o, errors.New("unable to get order from database")
	}

	deliveryQuery := fmt.Sprintf("SELECT %s from delivery where OrderUID = $1", deliveryFields)
	// Сбор данных об Delivery
	err = tx.QueryRow(deliveryQuery, oid).Scan(&o.Delivery.Name, &o.Delivery.Phone,
		&o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email)
	if err != nil {
		log.Printf("%v: unable to get delivery from database: %v\n", db.name, err)
		return o, errors.New("unable to get delivery from database")
	}

	paymentQuery := fmt.Sprintf(`SELECT %s FROM payment WHERE OrderUID = $1`, paymentFields)
	// Сбор данных о Payment
	err = tx.QueryRow(paymentQuery, oid).Scan(&o.Payment.Transaction,
		&o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank,
		&o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee)
	if err != nil {
		log.Printf("%v: unable to get payment from database: %v\n", db.name, err)
		return o, errors.New("unable to get payment from database")
	}

	itemQuery := fmt.Sprintf(`SELECT %s FROM item WHERE OrderUID = $1`, itemFields)
	// Сбор данных о Items
	rowsItems, err := tx.Query(itemQuery, oid)
	if err != nil {
		log.Printf("%v: unable to get items from database: %v\n", db.name, err)
		return o, errors.New("unable to get items from database")
	}
	defer rowsItems.Close()

	for rowsItems.Next() {
		var item Item
		err = rowsItems.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			log.Printf("%v: unable to scan item: %v\n", db.name, err)
			return o, errors.New("unable to scan item")
		}
		o.Items = append(o.Items, item)
	}

	return o, nil
}

// Сохранение Order в БД
func (db *DB) AddOrder(o Order) (int64, error) {
	tx, err := db.sqlxDB.Beginx()
	if err != nil {
		return 0, err
	}

	orderQuery := fmt.Sprintf(`INSERT INTO public.order (%s) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, orderFields)
	// Добавление Order
	_, err = tx.Exec(orderQuery, o.OrderUID, o.TrackNumber, o.Entry, o.Locale,
		o.InternalSignature, o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	deliveryQuery := fmt.Sprintf(`INSERT INTO delivery (OrderUID, %s)
		values ($1, $2, $3, $4, $5, $6, $7, $8)`, deliveryFields)
	// Добавление Delivery
	_, err = tx.Exec(deliveryQuery, o.OrderUID, o.Delivery.Name, o.Delivery.Phone,
		o.Delivery.Zip, o.Delivery.City, o.Delivery.Address, o.Delivery.Region, o.Delivery.Email)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	itemQuery := fmt.Sprintf(`INSERT INTO item (OrderUID, %s) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, itemFields)
	// добавление Items
	for _, item := range o.Items {
		_, err = tx.Exec(itemQuery, o.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Print(err)
			return 0, err
		}
	}

	paymentQuery := fmt.Sprintf(`INSERT INTO payment (OrderUID, %s) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, paymentFields)
	// Добавление Payment
	_, err = tx.Exec(paymentQuery, o.OrderUID, o.Payment.Transaction, o.Payment.RequestID,
		o.Payment.Currency, o.Payment.Provider, o.Payment.Amount, o.Payment.PaymentDt, o.Payment.Bank,
		o.Payment.DeliveryCost, o.Payment.GoodsTotal, o.Payment.CustomFee)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	log.Printf("%v: Order successfull added to DB\n", db.name)
	return 0, nil
}

func (db *DB) UploadCache() error {
	log.Println("Service is trying to upload cache...")
	tx, err := db.sqlxDB.Beginx()
	if err != nil {
		return err
	}

	var orders []Order
	getOrdersQuery := "SELECT * FROM public.order"
	getDeliveryQuery := fmt.Sprintf("SELECT %s FROM delivery WHERE OrderUID = $1", deliveryFields)
	getPaymentQuery := fmt.Sprintf("SELECT %s FROM payment WHERE OrderUID = $1", paymentFields)
	getItemsQuery := fmt.Sprintf("SELECT %s FROM item WHERE OrderUID = $1", itemFields)

	err = tx.Select(&orders, getOrdersQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, order := range orders {
		var delivery Delivery
		err = tx.Get(&delivery, getDeliveryQuery, order.OrderUID)

		if err != nil {
			tx.Rollback()
			return err
		}

		var payment Payment
		err = tx.Get(&payment, getPaymentQuery, order.OrderUID)

		if err != nil {
			tx.Rollback()
			return err
		}

		var items []Item
		err = tx.Select(&items, getItemsQuery, order.OrderUID)

		if err != nil {
			tx.Rollback()
			return err
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		db.cache.AddOrder(order)
	}

	log.Println("Upload of cache completed!")
	return tx.Commit()
}
