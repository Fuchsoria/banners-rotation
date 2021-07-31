# Banners rotation microservice ![](https://goreportcard.com/badge/github.com/Fuchsoria/banners-rotation)
API microservice with multi-armed bandit algorithm

## Commands
* **Start service containers:**
`make run`
* **Stop service containers:**
`make stop`
* **Start service in terminal with rebuild:**
`make up`
* **Lint service with golangci-lint:**
`make lint`
* **Run x100 unit tests with race:**
`make test`
* **Run x1 unit tests with race:**
`make t`
* **Run integration tests:**
`make integration-tests`
* **Regenerate grpc server and gateway:**
`make generate-gateway`

## Api endpoints
* **Create new banner, body:** `{"id":"","description":""}`
POST `/api/v1/admin/banners/create`
* **Create new slot, body:** `{"id":"","description":""}`
POST `/api/v1/admin/slots/create`
* **Create new social demo group, body:** `{"id":"","description":""}`
POST `/api/v1/admin/social-demos/create`
* **Add banner to rotation, body:** `{"banner_id":"","slot_id":""}`
POST `/api/v1/banners/add`
* **Add remove banner from rotation, body:** `{"banner_id":"","slot_id":""}`
POST `/api/v1/banners/remove`
* **Add click event, body:** `{"banner_id":"","slot_id":"","social_demo_id":""}`
POST `/api/v1/banners/click`
* **Get banner from slot, body:** `{"slot_id":"","social_demo_id":""}`
POST `/api/v1/banners/get`