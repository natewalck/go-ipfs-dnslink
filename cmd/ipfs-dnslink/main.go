package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/billputer/go-namecheap"
	"github.com/jpillora/go-tld"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:    "api_user",
				Usage:   "Enter the username for namecheap (if different from main account)",
				EnvVars: []string{"API_USER"},
			},
		),
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:    "api_token",
				Usage:   "api key for namecheap",
				EnvVars: []string{"API_TOKEN"},
			},
		),
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:    "user",
				Usage:   "Main namecheap account username",
				EnvVars: []string{"NAMECHEAP_USER"},
			},
		),
		altsrc.NewStringFlag(
			&cli.StringFlag{
				Name:    "domain",
				Usage:   "Domain to updated the dnslink entry on. Should be in the format of some.domain.com",
				EnvVars: []string{"NAMECHEAP_DOMAIN"},
			},
		),
		altsrc.NewIntFlag(
			&cli.IntFlag{
				Name:    "ttl",
				Value:   1800,
				EnvVars: []string{"NAMECHEAP_TTL"},
			},
		),
		&cli.StringFlag{Name: "cid", Usage: "IPFS CID to update the dnslink to."},
		&cli.StringFlag{Name: "config", Usage: "Path to config.yaml"},
	}

	app := &cli.App{
		Action: func(c *cli.Context) error {
			// Make sure we are fully configured before running
			for _, value := range []string{"api_user", "api_token", "user", "domain", "cid"} {
				if c.String(value) == "" {
					return fmt.Errorf("Missing value for %s", value)
				}
			}
			err := updateDNSLink(
				c.String("api_user"),
				c.String("api_token"),
				c.String("user"),
				c.String("domain"),
				c.String("cid"),
				c.Int("ttl"),
			)
			if err != nil {
				return err
			}
			return nil
		},
		Flags:       flags,
		Description: "App to quickly update dnslink records on namecheap",
		Version:     "0.0.1",
	}

	for _, argument := range os.Args {
		if strings.Contains(argument, "config") {
			app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func updateDNSLink(api_user, api_token, user, domain, cid string, ttl int) error {
	// Add http if missing, otherwise the TLD lib doens't parse correctly -_-
	var url string
	if strings.HasPrefix(domain, "https://") || strings.HasPrefix(domain, "http://") {
		url = domain
	} else {
		url = "http://" + domain
	}
	u, _ := tld.Parse(url)
	if u.Subdomain == "" {
		return fmt.Errorf("%s is not a valid full domain, please include the subdomain", domain)
	}
	client := namecheap.NewClient(api_user, api_token, user)
	// Re-assemble domain + TLD to compare to namecheap domains
	partial_domain := fmt.Sprintf("%s.%s", u.Domain, u.TLD)
	// Get a list of your domains
	namecheap_domains, _ := client.DomainsGetList()
	for _, namecheap_domain := range namecheap_domains {
		if namecheap_domain.Name == partial_domain {
			// We are only updating one record
			dnslink_record := []namecheap.DomainDNSHost{
				namecheap.DomainDNSHost{
					Name:    u.Subdomain,
					Type:    "TXT",
					Address: cid,
					TTL:     ttl,
				},
			}
			_, err := client.DomainDNSSetHosts(
				u.Domain,
				u.TLD,
				dnslink_record,
			)
			if err != nil {
				return err
			}
			fmt.Printf("Sucessfully updated %s to %s", domain, cid)
			return nil
		} else {
			continue
		}
	}
	return nil
}
