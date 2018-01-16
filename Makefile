OUTPUT_DIRECTORY = bin

build :
	mkdir -p $(OUTPUT_DIRECTORY)
	go get
	go build -o $(OUTPUT_DIRECTORY)/spectre