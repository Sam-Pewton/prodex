# Prodex - Productivity Index

A small lightweight application designed to scrape some basic metrics about work performance.
Designed to "lighten the load" of manually keeping track of all the work that you have done.


## Requirements
- systemd
- sqlite3


## Installation
Clone the repository and run the following from the project root:
```bash
make install
```
This will build the binary and install as a systemd service for the user (not root).
Two services are installed as a result:
- prodex.service (along with prodex.timer)
- prodexui.service

:warning: By default, the scraper will run on an hourly basis. If you would like to modify this
behaviour, you need to modify `./dev/prodex.timer` BEFORE running make install.

To uninstall, you can also run:
```bash
make uninstall
```

:warning: This will delete all configuration for the application and all scraped data. If you would
like to retain your configuration and data, make sure it is backed up somewhere before running. I
intend to fix this side effect later on once configurations are properly defined.


## Usage
### Scraper
After installation, the scraper is set to run hourly. There should be no need to modify this.
Any configuration changes made should be applied on the following run.

### UI
Currently, only a very basic web server is available on port `localhost:8642` that doesn't do a lot
of anything.

I will soon write a web UI that will make it easier to interact with the scraped data,
but in the meantime, you are able to query the database directly using a tool such as
nvim-dadbod or dbbrowser for sqlite3. The DB is located under `~/etc/prodex/db/prodex.db`.


## How it works
Two modes are available for prodex.
- Scraper
    This is the primary means of retrieving data and storing it in the DB.
    More TODO
- UI
    This is the primary means of viewing the data, but also provides a way to interact
    with user defined custom tables, for anything that needs recording but doesn't have
    an associated scraper (e.g. events organised).
    More TBD on how these tables will be created.

![architecture diagram](./images/prodex.png)


## Configuration
All scrapers are configured using the `prodex.toml` file. This file is located at the following
path `~/.config/prodex/prodex.toml` after installation.

Specific configurations are detailed under each available scraper below.


### Available scrapers
- Jira
    - Required Properties
        - atlassian_token
        - atlassian_user
        - atlassian_domain
    - Example
        ```toml
        [[scraper.jira]]
        atlassian_token = "abcdefg"
        atlassian_user = "your-email@email.email"
        atlassian_domain = "your-atlassian-url.atlassian.net"
        ```
- Confluence (not working currently)
    - Required Properties
        - atlassian_token
        - atlassian_user
        - atlassian_domain
    - Example
        ```toml
        [[scraper.confluence]]
        atlassian_token = "abcdefg"
        atlassian_user = "your-email@email.email"
        atlassian_domain = "your-atlassian-url.atlassian.net"
        ```

### Contributing
TODO
