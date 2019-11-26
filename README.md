# Start project:
## Start db
sudo docker run -p 27018:27017 --name cart_api_test -it -d mongo
## Run app
go run main.go
