build:
	mkdir -p ~/etc/prodex
	go build -o ~/etc/prodex prodex.go scraper.go ui.go

install: build
	if ! test -f ~/.config/prodex/prodex.toml; then mkdir -p ~/.config/prodex && cp dev/prodex.toml ~/.config/prodex/prodex.toml; fi
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
	rm -rf ~/etc/prodex/prodex
	rm -rf ~/.config/systemd/user/prodex*
