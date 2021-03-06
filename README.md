# PixivCaptor

PixivCaptor has a chrome extension to capture pixiv illusts. See [pixiv-captor-extension](https://github.com/aitsuki/pixiv-captor-extension).

## Build & Run

```bash
go build --tags json1 .
```

```bash
go run --tags json1 . -u username -p password

-P int Port (default 8080)
-db string Sqlite database path (default "./pixiv.db")
-log string Log file path (default "./pixiv.log")
-p string password
-u string username
```

## API

| Method | URI        | Description           | Params                                        |
| ------ | ---------- | --------------------- | --------------------------------------------- |
| HEAD   | /pixiv/:id | check illust captured | -                                             |
| GET    | /pixiv/:id | get illust            | `r18`: 0 - 1, `limit`: 1 ~ 100                |
| GET    | /pixiv     | search illust         | `r18`: 0 - 1, `limit`: 1 ~ 100, `q`: "string" |
| DELETE | /pixiv/:id | delete illust         | Basic Authorization                           |
| POST   | /pixiv     | capture illust        | Basic Authorization                           |

### POST - Capture

Request `www.pixiv.net/ajax/illust/:id` and `www.pixiv.net/ajax/illust/:id/pages`.

Then splicing pages body to illust body.

```text
{
    "id": "xxxxx", <--- illust
    "...": "...", <--- illust
    "pages":[]  <--- pages
}
```
