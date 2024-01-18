

NAME=redos
SRC=src/main.go


all: clean
	go build -o $(NAME) $(SRC)

scheck:
	staticcheck ./...

clean:
	rm -Rf $(NAME)

