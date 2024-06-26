<h2>Сервис для отображения данных о заказе.</h2>

<h3>Установка:</h3>

1. Клонировать репозиторий

2. Для запуска nats-streaming использовать команду

```bash
docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223
```

3. Для запуска сервера использовать команду

``` bash
go run cmd/main.go
```

В папке Schema добавлены миграции Postgres. Для корректного создания БД требуется:

1. Запустить postgres.

2. Создать БД С именем l0wb и в конфигах изменить настройки

3. Командой ниже мигрировать 

```bash
migrate -path ./schema -database 'postgres://DB_USER:DB_PASS@DB_HOST:DB_PORT/DB_NAME?sslmode=disable' up
```

4. Если вы хотите отменить миграцию воспользуйтесь командой

```bash
migrate -path ./schema -database 'postgres://DB_USER:DB_PASS@DB_HOST:DB_PORT/DB_NAME?sslmode=disable' down
```

<h4>Чтобы убедиться в работе сервиса необходимо зайти на http://localhost:3333/

По этому адресу отображаются все существующие Order.

Указав правильный OrderUID и нажав Search вам отобразится конкретный Order</h4>
