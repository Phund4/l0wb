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
migrate -path ./schema -database 'postgres://phunda:098908@localhost:5432/l0wb?sslmode=disable' up
```

4. Если вы хотите отменить миграцию воспользуйтесь командой

```bash
migrate -path ./schema -database 'postgres://phunda:098908@localhost:5432/l0wb?sslmode=disable' down
```