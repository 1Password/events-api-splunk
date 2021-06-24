.PHONY: compile_app_binary
compile_app_binary:
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=signinattempts" -o ../app/op_events_reporting/bin/signin_attempts
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=itemusages" -o ../app/op_events_reporting/bin/item_usages
	@cp -R src app/op_events_reporting/lib/item_usages
	@cp -R src app/op_events_reporting/lib/signin_attempts

.PHONY: clean
clean:
	@cd app && $(MAKE) clean
.PHONY: build_all_binaries
build_all_binaries:
	@rm -rf builds/bin && mkdir builds/bin
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/op_events_reporting/bin/signin_attempts" -ldflags '-s -X main.EventBuildType=signinattempts' .
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/op_events_reporting/bin/item_usages" -ldflags '-s -X main.EventBuildType=itemusages' .

.PHONY: build_all_apps
build_all_apps: clean
	@cp -R src app/op_events_reporting/lib/item_usages
	@cp -R src app/op_events_reporting/lib/signin_attempts
	@cd builds/bin && for d in */; do cp -a ../../app/op_events_reporting $${d}; done
	@sed -i'.bak' 's#bin/signin_attempts#bin/signin_attempts.exe#g' builds/bin/windows_amd64/op_events_reporting/default/inputs.conf
	@sed -i'.bak' 's#bin/item_usages#bin/item_usages.exe#g' builds/bin/windows_amd64/op_events_reporting/default/inputs.conf
	@rm -f builds/bin/windows_amd64/op_events_reporting/default/inputs.conf.bak
	@cd builds/bin && for d in */; do \
		cd $${d}; \
		COPYFILE_DISABLE=1 tar --exclude='.DS_Store' --exclude='.gitignore' --exclude='.travis.yml' -cvzf op_events_reporting.tar.gz op_events_reporting; \
		cd ..; \
	done
	@cd builds/bin && for d in */; do rm -rf $${d}op_events_reporting; done
