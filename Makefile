build:
	mkdir -p ~/etc/prodex
	go build -o ~/etc/prodex prodex.go scraper.go ui.go

install: build
	if ! test -f ./.env; then if ! test -f ~/.config/prodex/.env; then echo "No .env file located in cwd. Exiting" exit 1; fi; fi
	mkdir -p ~/.config/prodex
	cp .env ~/.config/prodex/.env
	mkdir -p ~/.config/systemd/user
	cp dev/prodex.service ~/.config/systemd/user/prodex.service
	cp dev/prodex.timer ~/.config/systemd/user/prodex.timer
	cp dev/prodexui.service ~/.config/systemd/user/prodexui.service
	systemctl --user enable prodex
	systemctl --user enable prodex.timer
	systemctl --user start prodex
	systemctl --user enable prodexui
	systemctl --user start prodexui

uninstall:
	systemctl --user stop prodexui
	systemctl --user disable prodexui
	systemctl --user stop prodex
	systemctl --user disable prodex
	systemctl --user disable prodex.timer
	rm -rf ~/etc/prodex
	rm -rf ~/.config/prodex
	rm -rf ~/.config/systemd/user/prodex*
