sleep 5 && git pull && if pgrep ResServiceBot; then killall ResServiceBot; fi && go build && ./ResServiceBot