# Banner Service
сервис для управления разнородным контентом

## Installation & Run
```bash
# Скачать проект
git clone https://github.com/CyberPiess/banner_service
```
Перед запуском API сервера необходимо настроить переменные окружения. Для этого нужно создать файл .env в директории /build по аналогии с файлом .env.local 
```.env
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
PG_HOST= `pg_image_host` #должен соответствовать имени контейнера с postgres
PG_PORT=5432
SSLMODE=disable
DBNAME=`db_name` default: banner_service
REDIS_PASSWORD=somepassword
REDIS_ADRESS=redis:6379
```
Компиляция и запуск осуществляются при помощи Makefile

* make выполняет запуск контейнера с postgres, создание базы данных, миграцию таблиц, запуск контейнера Redis, сборку и запуск API
- make install отвечает за установку библиотеки для осуществления миграций
- make dropdb удаляет базу данных
- make migratedown откатывает изменения внесенные миграцией

```bash
# Компиляция и запуск
cd banner_service
make install # если нет библиотеки для осуществления миграций
make
```

## API

#### /api/user_banner?tag_id=1&feature_id=1
- *Method* : `GET`
    - *Parameters* :
      - `tag_id` (query, required):идентификатор тега
      - `feature_id` (query, required): идентификатор фичи
      - `use_last_revision` (query, default: false): признак получения наиболее актуальной информации
      - `token` (header, required): токен пользователя или админа
    - *Description*: Получение баннера для пользователя
  
#### /api/banner
1. *Method* : `GET`
    - *Parameters* :
      - `tag_id` (query, optional): идентификатор тега
      - `feature_id` (query, optional): идентификатор фичи
      - `limit` (query, optional): лимит
      - `offset` (query, optional): оффсет
      - `token` (header, required): токен админа
    - *Description*:  Получение всех баннеров c фильтрацией по фиче и/или тегу
2. *Method* : `POST`
    - *Parameters* :
      - `tag_ids` (body, required): идентификаторы тэгов
      - `feature_id` (body, required): идентификатор фичи
      - `content` (body, required): cодержимое баннера
      - `is_active` (body, required): флаг активности баннера
      - `token` (header, required): токен админа
    - *Description*:  Создание нового баннера

#### /api/banner/{id}
1. *Method* : `PUT`
    - *Parameters* :
      - `id` (path, required): идентификатор баннера
      - `tag_ids` (body, required, nullable): идентификаторы тэгов
      - `feature_id` (body, required, nullable): идентификатор фичи
      - `content` (body, required, nullable): cодержимое баннера
      - `is_active` (body, required, nullable): флаг активности баннера
      - `token` (header, required): токен админа
    - *Description*:  Обновление содержимого баннера
2. *Method* : `DELETE`
    - *Parameters* :
      - `id` (path, required): идентификатор баннера
      - `token` (header, required): токен админа
    - *Description*:  Удаление баннера по идентификатору

  ----
Примеры запросов и ответов находятся в директории postman_collections, файл banner_service_http_requests.json
