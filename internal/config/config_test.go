package config

import (
	"fmt"
	"testing"
)

// define the required vars, and the three scrapers.
// 2 scrapers for jira and 1 for confluence.
func TestLoadGoodConfig(t *testing.T) {
	config := `
    installation_path = "/etc/prodex/"
    log_level = "debug"
    max_noops = 10000
    
    [scrapers]
    [[scrapers.jira]]
    atlassian_token = "abcd"
    atlassian_user = "test"
    atlassian_domain = "https://test.atlassian.net"
    pagination_size = 50

    [[scrapers.jira]]
    atlassian_token = "abcd"
    atlassian_user = "test"
    atlassian_domain = "https://test.atlassian.net"
    pagination_size = 50

    [[scrapers.confluence]]
    atlassian_token = "abcd"
    atlassian_user = "test"
    atlassian_domain = "https://test.atlassian.net"
    pagination_size = 50
    `

	err := LoadConfig(config)

	if err != nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}

	if ProdexConf.LogLevel != "debug" {
		t.Fatalf(fmt.Sprintf("got: %s", ProdexConf.LogLevel))
	}
	if ProdexConf.MaxNoops != 10000 {
		t.Fatalf("max noops was not correct")
	}
	if len(ProdexConf.Scrapers) != 2 {
		t.Fatalf("not enough scrapers defined")
	}
	if len(ProdexConf.Scrapers["jira"]) != 2 {
		t.Fatalf("not enough scrapers defined")
	}
	if len(ProdexConf.Scrapers["confluence"]) != 1 {
		t.Fatalf("not enough scrapers defined")
	}
}

// Expect an error, as the variable has no assignment in the toml
func TestLoadBadConfig(t *testing.T) {
	config := `max_noops`

	err := LoadConfig(config)

	if err == nil {
		t.Fatalf(fmt.Sprintf("%s", err))
	}
}
