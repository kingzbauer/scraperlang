scanner:
	go build -o scan github.com/kingzbauer/scraperlang/cmd/scanner

make-parser:
	go build -o parse github.com/kingzbauer/scraperlang/cmd/parser

clean:
	-rm parse scan
