NAME = train
NAME2 = predict
SRC = train.go
SRC2 = predict.go

GREEN = \033[0;32m
RED = \033[0;31m
NC = \033[0m

all: $(NAME)

$(NAME):
	@mkdir -p img 
	@go mod init $(NAME)
	@go get gonum.org/v1/plot/...
	@go get gonum.org/v1/plot/plotter
	@go get gonum.org/v1/plot/vg
	@go build -o $(NAME) $(SRC)
	@go build -o $(NAME2) $(SRC2)

clean:
	@rm -rf $(NAME) $(NAME2) go.mod go.sum model.json img model.txt

re: clean all

.PHONY: all clean re
