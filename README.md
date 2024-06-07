# spender 
Spender - is opensource api for register spending of money and time

## Project outline
- cmd -> command directory
  - api -> api server main package
- handlers -> api handle routes
- types -> api models
- db -> store api and realisation
- utils -> usefull functions
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

### Air
Documentation
```
https://github.com/cosmtrek/air/blob/master/README.md
```

Installing air framework
```
go install github.com/cosmtrek/air@latest
```

## Docker
### Load mongodb as docker container
pull and run container
```
docker run --name spender-mongo -d mongo:latest -p 27017:27017
```

run container
```
docker start spender-mongo
```