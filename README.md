<h2>Сервис для отображения данных о заказе.</h2>

<h3>Установка:</h3>
1. Клонировать репозиторий
<br/>
2. Для запуска nats-streaming использовать команду
<br/>

```bash
docker run -p 4223:4223 -p 8223:8223 nats-streaming -p 4223 -m 8223
```

3. Для запуска сервера использовать команду
<br/>

``` bash
go run cmd/main.go
```

В папке Schema добавлены миграции Postgres. Для корректного создания БД требуется:
<br/>
1. Запустить postgres.
<br/>
2. Создать БД С именем l0wb и в конфигах изменить настройки
<br/>
3. Командой ниже мигрировать 
<br/>

```bash
migrate -path ./schema -database 'postgres://phunda:098908@localhost:5432/l0wb?sslmode=disable' up
```

4. Если вы хотите отменить миграцию воспользуйтесь командой
<br/>

```bash
migrate -path ./schema -database 'postgres://phunda:098908@localhost:5432/l0wb?sslmode=disable' down
```