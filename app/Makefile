# text reset
NO_COLOR=\033[0m
# green
OK_COLOR=\033[32;01m
# red
ERROR_COLOR=\033[31;01m
# cyan
WARN_COLOR=\033[36;01m
# yellow
ATTN_COLOR=\033[33;01m

ROOT_DIR := $(shell git rev-parse --show-toplevel)

VERSION := `git describe --tags --dirty 2>/dev/null`
COMMITHASH := `git rev-parse --short HEAD 2>/dev/null`
DATE := `date "+%FT%T%z"`

CONTAINER_NAME := 'op_events_reporting'

.PHONY: all
all: restart

init:
	@echo "$(ATTN_COLOR)==> init $(NO_COLOR)"

.PHONY: up
up:
	@cd .. && $(MAKE) compile_app_binary
	@touch splunkd.log
	@echo "$(ATTN_COLOR)==> up $(NO_COLOR)"
	@docker-compose up -d

.PHONY: remove
remove:
	@echo "$(ATTN_COLOR)==> rm $(NO_COLOR)"
	@docker-compose rm -f -s

.PHONY: wait_up
wait_up:
	@echo "$(ATTN_COLOR)==> wait_up $(NO_COLOR)"
	@for i in `seq 0 180`; do if docker exec -it $(CONTAINER_NAME) /sbin/checkstate.sh &> /dev/null; then break; fi; printf "\rWaiting for Splunk for %s seconds..." $$i; sleep 1; done

.PHONY: down
down:
	@echo "$(ATTN_COLOR)==> down $(NO_COLOR)"
	@docker-compose stop

.PHONY: start
start: up wait_up

.PHONY: restart
restart: down start

.PHONY: refresh
refresh: remove start

.PHONY: clean
clean:
	$(RM) -rf op_events_reporting/local
	$(RM) op_events_reporting/metadata/local.meta
	$(RM) op_events_reporting/bin/item_usages
	$(RM) -rf op_events_reporting/lib/item_usages
	$(RM) op_events_reporting/bin/signin_attempts
	$(RM) -rf op_events_reporting/lib/signin_attempts
	$(RM) splunkd.log
	@touch splunkd.log
