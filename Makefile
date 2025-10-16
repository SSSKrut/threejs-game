.PHONY: backend-dev
backend-dev:
	cd backend && go run .

.PHONY: go-mod-tidy
go-mod-tidy:
	cd backend && go mod tidy

.PHONY: frontend-dev
frontend-dev:
	cd game && yarn dev

.PHONY: frontend-build
frontend-build:
	cd game && yarn build

.PHONY: frontend-preview
frontend-preview:
	cd game && yarn preview

.PHONY: yarn-install
yarn-install:
	cd game && yarn install

.PHONY: dev
dev: backend-dev frontend-dev

.PHONY: build
build: frontend-build