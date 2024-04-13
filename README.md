# Banner Service
сервис для управления разнородным контентом

## Installation & Run

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







