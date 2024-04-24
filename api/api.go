package api

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/Phund4/l0wb/internal/db"

	"github.com/go-chi/chi/v5"
)

type ordkey string

const orderKey ordkey = "order"

type Api struct {
	router                *chi.Mux
	cache              *db.Cache
	name               string
	server                *http.Server
	httpServerExitDone *sync.WaitGroup
}

func NewApi(csh *db.Cache) *Api {
	api := Api{}
	api.Init(csh)
	return &api
}

// Инициализация сервера
func (a *Api) Init(csh *db.Cache) {
	a.cache = csh
	a.name = "API"
	a.router = chi.NewRouter()
	a.router.Get("/", a.WelcomeHandler)

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

// Корректное завершение работы сервера
func (a *Api) Finish() {
	if err := a.server.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	a.httpServerExitDone.Wait()
	log.Printf("%v: Сервер успешно выключен!\n", a.name)
}

// Запуск сервера в отдельном потоке (для корректного завершения работы программы: очистка кеша из БД, отключение от подписки)
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
	t, err := template.ParseFiles("UI/order.html")
	if err != nil {
		log.Printf("%v: getOrder(): ошибка парсинга шаблона html: %s\n", a.name, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = t.ExecuteTemplate(w, "order.html", nil)
	if err != nil {
		log.Printf("%v: WellcomeHandler(): ошибка выполнения шаблона html: %s\n", a.name, err)
		return
	}
}

// Мидлвара, сохраняющая в контекст Order
func (a *Api) orderCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orderIDstr := chi.URLParam(r, "orderID")
		orderID, err := strconv.ParseInt(orderIDstr, 10, 64)
		if err != nil {
			log.Printf("%v: ошибка конвертации %s в число: %v\n", a.name, orderIDstr, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		log.Printf("%v: запрос OrderOut из кеша/бд, OrderID: %v\n", a.name, orderIDstr)
		order, err := a.cache.GetOrder(fmt.Sprint(orderID))
		orderOut := db.OrderOut{OrderUID: order.OrderUID, TrackNumber: order.TrackNumber, Entry: order.Entry}
		if err != nil {
			log.Printf("%v: ошибка получения OrderOut из кеша/бд: %v\n", a.name, err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), orderKey, orderOut)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Хендлер запроса Order
func (a *Api) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderOut, ok := ctx.Value(orderKey).(db.OrderOut)
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
	err = t.ExecuteTemplate(w, "order.html", orderOut)
	if err != nil {
		log.Printf("%v: GetOrder(): ошибка выполнения шаблона html: %s\n", a.name, err)
		return
	}
}
