# We only want to generate the parser when there's changes to the grammar.
.PHONY: parser
parser:
	@echo ">> running goyacc to generate the .go file."
ifeq (, $(shell command -v goyacc > /dev/null))
	@echo "goyacc not installed so skipping"
	@echo "To install: go install golang.org/x/tools/cmd/goyacc@v0.6.0"
else
	goyacc -o pkg/b_promql/parser/generated_parser.y.go pkg/b_promql/parser/generated_parser.y
endif