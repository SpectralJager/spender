# spender 
Spender - is opensource api for register spending of money and time

## Resources
### Mongodb driver
Documentation
```
https://www.mongodb.com/docs/drivers/go/current/quick-start/
```

Installing mongodb driver
```
go get go.mongodb.org/mongo-driver/mongo
```

### Echo v4
Documentation
```
https://echo.labstack.com/docs
```

Installing echo framework
```
go get github.com/labstack/echo/v4
```

## Docker
### Load mongodb as docker container
```
docker run --name spender-mongo -d mongo:latest -p 27017:27017
```