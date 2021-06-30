.PHONY: compile_app_binary
compile_app_binary:
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=signinattempts" -o ../app/onepassword_events_api/bin/signin_attempts
	@cd src && env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.EventBuildType=itemusages" -o ../app/onepassword_events_api/bin/item_usages
	@cp -R src app/onepassword_events_api/lib/item_usages
	@cp -R src app/onepassword_events_api/lib/signin_attempts

.PHONY: clean
clean:
	@cd app && $(MAKE) clean
.PHONY: build_all_binaries
build_all_binaries:
	@rm -rf builds/bin && mkdir builds/bin
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/onepassword_events_api/bin/signin_attempts" -ldflags '-s -X main.EventBuildType=signinattempts' .
	@cd src && gox -arch="amd64 arm" -os="linux windows freebsd openbsd" -osarch="darwin/amd64" -output="../builds/bin/{{.OS}}_{{.Arch}}/onepassword_events_api/bin/item_usages" -ldflags '-s -X main.EventBuildType=itemusages' .

.PHONY: build_all_apps
build_all_apps: clean
	@cp -R src app/onepassword_events_api/lib/item_usages
	@cp -R src app/onepassword_events_api/lib/signin_attempts
	@cd builds/bin && for d in */; do cp -a ../../app/onepassword_events_api $${d}; done
	@sed -i'.bak' 's#bin/signin_attempts#bin/signin_attempts.exe#g' builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf
	@sed -i'.bak' 's#bin/item_usages#bin/item_usages.exe#g' builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf
	@rm -f builds/bin/windows_amd64/onepassword_events_api/default/inputs.conf.bak
	@cd builds/bin && for d in */; do \
		cd $${d}; \
		COPYFILE_DISABLE=1 tar --exclude='.DS_Store' --exclude='.gitignore' --exclude='.travis.yml' --exclude='.gitcookies.sh.enc' -cvzf onepassword_events_api.tar.gz onepassword_events_api; \
		cd ..; \
	done
	@cd builds/bin && for d in */; do rm -rf $${d}onepassword_events_api; done
