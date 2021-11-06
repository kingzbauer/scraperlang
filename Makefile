scanner:
	go build -o scan github.com/kingzbauer/scraperlang/cmd/scanner

make-parser:
	go build -o parse github.com/kingzbauer/scraperlang/cmd/parser

interpret:
	go build -o sl github.com/kingzbauer/scraperlang/cmd/interpreter

clean:
	-rm parse scan sl
