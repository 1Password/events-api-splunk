# Run make new_version after changing this version
VERSION=1.14.0

.PHONY: compile_app_binary
compile_app_binary:
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=signinattempts -X go.1password.io/eventsapi-splunk/api.Version=$(VERSION)" -o ../app/onepassword_events_api/bin/signin_attempts
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=itemusages -X go.1password.io/eventsapi-splunk/api.Version=$(VERSION)" -o ../app/onepassword_events_api/bin/item_usages
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=auditevents -X go.1password.io/eventsapi-splunk/api.Version=$(VERSION)" -o ../app/onepassword_events_api/bin/audit_events
	@cp -R src app/onepassword_events_api/lib/item_usages
	@cp -R src app/onepassword_events_api/lib/signin_attempts
	@cp -R src app/onepassword_events_api/lib/audit_events

.PHONY: clean
clean:
	@cd app && $(MAKE) clean
.PHONY: build_all_binaries
build_all_binaries:
	@rm -rf builds/bin && mkdir builds/bin
	@cd app/onepassword_events_api && npm run build-release
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/onepassword_events_api/bin/signin_attempts" -ldflags '-s -X main.EventBuildType=signinattempts -X go.1password.io/eventsapi-splunk/api.Version=$(VERSION)' .
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/onepassword_events_api/bin/item_usages" -ldflags '-s -X main.EventBuildType=itemusages -X go.1password.io/eventsapi-splunk/api.Version=$(VERSION)' .
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/onepassword_events_api/bin/audit_events" -ldflags '-s -X main.EventBuildType=auditevents -X go.1password.io/eventsapi-splunk/api.Version=$(VERSION)' .

.PHONY: build_all_apps
build_all_apps: clean
	@cp -R src app/onepassword_events_api/lib/item_usages
	@cp -R src app/onepassword_events_api/lib/signin_attempts
	@cp -R src app/onepassword_events_api/lib/audit_events
	@cd builds/bin && for d in */; do cp -a ../../app/onepassword_events_api $${d}; done
	@sed -i'.bak' 's#bin/signin_attempts#bin/signin_attempts.exe#g' builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf
	@sed -i'.bak' 's#bin/item_usages#bin/item_usages.exe#g' builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf
	@sed -i'.bak' 's#bin/audit_events#bin/audit_events.exe#g' builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf
	@rm -f builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf.bak
	@cd builds/bin && for d in */; do \
		cd $${d}; \
		COPYFILE_DISABLE=1 tar --exclude='.DS_Store' --exclude='.gitignore' --exclude='.travis.yml' --exclude='.gitcookies.sh.enc' --exclude="node_modules" --exclude="package.json" --exclude="package-lock.json" --exclude="webpack.config.js" -cvzf onepassword_events_api_$(VERSION).tar.gz onepassword_events_api; \
		cd ..; \
	done
	@cd builds/bin && for d in */; do rm -rf $${d}onepassword_events_api; done

.PHONY: new_version
new_version:
	@cd app/onepassword_events_api && npm version $(VERSION) && node util/build_version.js && npm run build-release
