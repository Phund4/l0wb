package api

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/Phund4/l0wb/internal/db"

	"github.com/go-chi/chi/v5"
)

// Ключи для context
type orderkey string;
type orderskey string;
const orderKey orderkey = "order"
const ordersKey orderskey = "orders"

// Структура Api
type Api struct {
	router             *chi.Mux
	cache              *db.Cache
	name               string
	server             *http.Server
	httpServerExitDone *sync.WaitGroup
}

// Создание Api
func NewApi(csh *db.Cache) *Api {
	api := Api{}
	api.Init(csh)
	return &api
}

// Инициализация Api
func (a *Api) Init(csh *db.Cache) {
	a.cache = csh
	a.name = "API"
	a.router = chi.NewRouter()

	a.router.Route("/", func(r chi.Router) {
		r.Use(a.ordersCtx)
		r.Get("/", a.WelcomeHandler)
	})
	
	a.router.Route("/orders", func(r chi.Router) {
		r.Route("/{orderID}", func(r chi.Router) {
			r.Use(a.orderCtx)
			r.Get("/", a.GetOrder)
		})
	})

	a.httpServerExitDone = &sync.WaitGroup{}
	a.httpServerExitDone.Add(1)
	a.StartServer()
}

// Завершение работы сервера
func (a *Api) Finish() {
	if err := a.server.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	a.httpServerExitDone.Wait()
	log.Printf("%v: Сервер успешно выключен!\n", a.name)
}

// Запуск сервера в отдельном потоке
func (a *Api) StartServer() {
	a.server = &http.Server{
		Addr:    ":3333",
		Handler: a.router,
	}

	go func() {
		defer a.httpServerExitDone.Done()

		log.Printf("%v: сервер будет запущен по адресу http://localhost:3333\n", a.name)
		if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("ListenAndServe() error: %v", err)
			return
		}
	}()
}

// Обработчик главной страницы http://localhost:3333
func (a *Api) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orders, ok := ctx.Value(ordersKey).([]db.Order)
	if !ok {
		log.Printf("%v: WellcomeHandler(): ошибка приведения интерфейса к типу []Order\n", a.name)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	t, err := template.ParseFiles("UI/orders.html")
	if err != nil {
		log.Printf("%v: WellcomeHandler(): ошибка парсинга шаблона html: %s\n", a.name, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	type ords struct {
		Orders []db.Order
	}

	err = t.ExecuteTemplate(w, "orders.html", ords{Orders: orders})
	if err != nil {
		log.Printf("%v: WellcomeHandler(): ошибка выполнения шаблона html: %s\n", a.name, err)
		return
	}
}

// Middleware, сохраняющий в контекст Orders
func (a *Api) ordersCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v: запрос Orders из кеша/бд\n", a.name)
		orders, err := a.cache.GetOrders()
		if err != nil {
			log.Printf("%v: ошибка получения Orders из кеша/бд: %v\n", a.name, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ordersKey, orders)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware, сохраняющий в контекст Order
func (a *Api) orderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orderIDstr := chi.URLParam(r, "orderID")

		log.Printf("%v: запрос OrderOut из кеша/бд, OrderID: %v\n", a.name, orderIDstr)
		order, err := a.cache.GetOrder(orderIDstr)
		if err != nil {
			log.Printf("%v: ошибка получения OrderOut из кеша/бд: %v\n", a.name, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), orderKey, order)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Хендлер запроса Order
func (a *Api) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	order, ok := ctx.Value(orderKey).(db.Order)
	if !ok {
		log.Printf("%v: GetOrder(): ошибка приведения интерфейса к типу OrderOut\n", a.name)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	t, err := template.ParseFiles("UI/order.html")
	if err != nil {
		log.Printf("%v: GetOrder(): ошибка парсинга шаблона html: %s\n", a.name, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = t.ExecuteTemplate(w, "order.html", order)
	if err != nil {
		log.Printf("%v: GetOrder(): ошибка выполнения шаблона html: %s\n", a.name, err)
		return
	}
}
