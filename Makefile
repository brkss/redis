

NAME=redos
SRC=src/main.go


all: clean
	go build -o $(NAME) $(SRC)

clean:
	rm -Rf $(NAME)

