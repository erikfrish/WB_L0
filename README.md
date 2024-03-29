# WB_L0

### Как запустить

Чтобы запустить приложение, достаточно скопировать репозиторий и поднять контейнеры через docker compose. Выполните в корне проектра команду:

```bash
docker compose up
```

При этом запустится основное приложение sub, которое является Subscriber'ом в nats-streaming, а также веб сервером, которвый выдает данные о заказах по order_uid.

Чтобы получить данные о заказе из кэша, можно перейти на простую веб-страницу запущенную в контенере nginx на localhost:80.

Чтобы передать новые данные в sub через nats-streaming можно запустить приложение cmd/pub/main.go, оно возьмет моковые данные из папки mock_orders и начнет посылать в nats-streaming с интервалом в 1 секунду. В коде подключен один из файлов с моковыми данными, в каждом файле 1000 заказов. Запустить pub можно локально на хостовой машине, достаточно набрать в терминале:

```bash
go run cmd/pub/main.go
```

### Список дел:

- [X] Сделать кэш
- [X] Сделать БД
- [X] Подключение к БД, внесение строк, получение строк по одной и все сразу

#### СР

- [X] Сделать веб сервер в сабе, который будет обрабатывать GET запросы на получение информации по заказу по его айди с помощью chi
- [X] Сделать веб-морду для отправки запроса и отображения результата

#### ЧТ

- [X] Написать взаимодействие со СТАН для паба и саба
- [X] Реализовать запись полученных сабом данных в кэш и в бд (сначала в бд)

#### ПТ

- [X] Сделать валидацию данных при получении новых данных сабом, поможет github.com/go-playground/validator/v10
- [X] Сделать генератор данных по заказу (Магазин)
- [X] Подретушировать логи во всех пакетах под слог

#### СБ

- [ ] Написать тесты (со \*)

#### ВС

- [ ] Сдать готовый проект
