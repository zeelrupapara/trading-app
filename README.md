# Trading App

### Local Setup guide
#### 1) Start Database
```
docker compose up -d
```

#### 2) Migrations for DB
```
cd api 
cp .env.example .env
go mod tidy
go run main.go migrate up
```

#### 3) Start the Backend
```
go run main.go api
```

now backend running in localhost:8000

see swagger docs: http://localhost:8000/api/v1/docs


#### 4) Start the frontend
```
cd web 
npm install
npm run dev
```

now fronend running in http://localhost:5173

for see the live trade the best option is ETHBTC
If you want to see the live trade other then given then try to change route in trade page
(ex: http://localhost:5173/trades?symbol={symbol}) (ex: symbol: ETHBTC)