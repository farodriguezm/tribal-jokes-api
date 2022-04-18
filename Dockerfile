FROM golang

COPY . .

RUN go build main.go

CMD ./main
