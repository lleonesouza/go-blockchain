FROM golang:alpine

WORKDIR /app/

COPY . .

RUN go get github.com/githubnemo/CompileDaemon

EXPOSE 5000

CMD CompileDaemon --directory="." --build="go build ./main.go" --command=./main
