# puredns

[massdns](https://github.com/blechschmidt/massdns) is an incredibly powerful DNS resolver used to perform bulk lookups. With the proper bandwidth and a good list of public resolvers, it can resolve millions of queries in just a few minutes.

Unfortunately, the results are only as good as the answers provided by the public resolvers. They are often polluted by wildcards or DNS poisoned entries.

**puredns** is a bash application that uses massdns to _accurately_ perform DNS bruteforcing and mass resolving. It ensures the results obtained by public resolvers are clean and can work around DNS poisoning by validating the answers obtained using a list of trusted resolvers. It also handles wildcard subdomains through its companion Python script, **wildcarder**, which can also be used by itself.

## Features

* Accurately resolves thousands of DNS queries per second using massdns and a list of public DNS resolvers
* Supports DNS bruteforcing using a wordlist and a root domain
* Cleans wildcard root subdomains with a minimal number of queries
* Validates that the results are free of DNS poisoning by running against a list of known, trustable resolvers
* Saves clean massdns output, wildcard root domains and answers, and list of valid subdomains
* Comes with a massive list of tested public resolvers from [public-dns.info](https://public-dns.info/)

## Installation

puredns requires massdns to be installed on the host machine. [Follow the instructions](https://github.com/blechschmidt/massdns#compilation) to compile massdns on your system. It also needs to be in accessible through the PATH environment variable.

The script also requires a few other dependencies:

* python3
* pv

To ensure the dependencies are installed on Ubuntu:

```
sudo apt install -y python3 pv
```

Once the dependencies are installed, simply clone the repository to start using puredns:

```
git clone https://github.com/d3mondev/puredns.git
```

## Usage

### puredns

puredns perfoms three steps automatically:

#### 1. Dumb mass resolve

First, it performs a mass resolve of all the domains using massdns and saves the results to an intermediary file.

The results are usually polluted: some public resolvers will have sent poisoned replies, and wildcard subdomains can also quickly inflate the results.

This step is mandatory and cannot be skipped.

#### 2. Wildcard detection

The application then uses its companion _wildcarder_ script to detect and extract all the wildcard root subdomains from the massdns results file.

wildcarder will reuse the massdns output to minimize the number of queries it needs to perform. It is quite common for wildcarder to figure out the structure of wildcard subdomains with as few as 5-10 DNS queries when it is passed a massdns file to prime its cache.

puredns will then clean up all the wildcard subdomains from the results.

This step can be skipped using the `--skip-wildcard-check` flag.

#### 3. Validation

To protect against DNS poisoning, massdns is used one last time to validate the remaining results using a list of trusted DNS resolvers.

This step is done at a slower pace as not to hit any rate limiting on the trusted resolvers.

This step can be skipped using the `--skip-validation` flag.

Once this is done, the resulting files are clean of wildcard subdomains and DNS poisoned answers.

#### Command line arguments

```
puredns v1.0
Use massdns to accurately resolve a large amount of subdomains and extract wildcard domains.

Usage:
        puredns [--skip-validation] [--skip-wildcard-check] [--write-massdns <filename>]
                [--write-wildcards <filename>] [--write-wildcard-answers <filename>] [--help] <command> <args>

        Example:
                puredns [args] resolve domains.txt
                puredns [args] bruteforce wordlist.txt domain.com

        Commands:
                resolve <filename>              Resolve a list of domains
                bruteforce <wordlist> <domain>  Perform subdomain bruteforcing on a domain using a wordlist

        Optional:
                -r,  --resolvers <filename>             Text file containing resolvers
                -tr, --trusted-resolvers <filename>     Text file containing trusted resolvers

                -ss, --skip-sanitize                    Do not sanitize the list of domains to test
                                                        By default, domains are set to lowercase and
                                                        only valid characters are kept
                -sv, --skip-validation                  Do not validate massdns results using trusted resolvers
                -sw, --skip-wildcard-check              Do no perform wildcard detection and filtering

                -w, --write <filename>                  Write valid domains to a file
                -wm, --write-massdns <filename>         Write massdns results to a file
                -ww, --write-wildcards <filename>       Write wildcard root subdomains to a file
                -wa, --write-answers <filename>         Write wildcard DNS answers to a file

                -h, --help                              Display this message
```

#### Examples

Using the included list of resolvers, puredns can bruteforce a massive list of subdomains using a wordlist named _all.txt_:

`puredns bruteforce all.txt domain.com`

Only the resulting valid domains are sent to stdout so that you can pipe the results to other programs.

It can also save all the DNS answers and wildcard subdomains found:

`puredns resolve domains.txt --write-massdns dns_answers.txt --write-wildcards wildcards.txt --write valid_domains.txt`

You can specify your own custom resolvers:

`puredns resolve domains.txt -r resolvers.txt`

### wildcarder

wildcarder can extract wildcard root subdomains and their DNS answers from a list of domains.

It is used as part of the puredns workflow, but can also be used by itself.

#### Command line arguments

```
usage: wildcarder [-h] [--load-massdns-cache massdns.txt]
                  [--write-domains domains.txt] [--write-answers answers.txt]
                  [--version]
                  file

Find wildcards from a list of subdomains, optionally loading a massdns simple
text output file to reduce number of DNS queries.

positional arguments:
  file                  file containing list of subdomains

optional arguments:
  -h, --help            show this help message and exit
  --load-massdns-cache massdns.txt
                        load a DNS cache from a massdns output file (-oS)
  --write-domains domains.txt
                        write wildcard domains to file
  --write-answers answers.txt
                        write wildcard DNS answers to file
  --version             show program's version number and exit
```

## Resources

[shuffleDNS](https://github.com/projectdiscovery/shuffledns) is a good alternative written in go that handles wildcard subdomains using a different algorithm.

[public-dns.info](https://public-dns.info/) continuously updates a list of public and free DNS resolvers.

[DNS Validator](https://github.com/vortexau/dnsvalidator) can be used to curate your own list of public DNS resolvers.

## Author

I'm d3mondev, a coder and security enthusiast. [Follow me on twitter](https://twitter.com/d3mondev)!

## Disclaimer & License

The resolvers included in this repository are for reference only. The author is not responsible for any misuse of the resolvers in that list.

Usage of this program for attacking targets without consent is illegal. It is the user's responsibility to obey all applicable laws. The developer assumes no liability and is not responsible for any misuse or damage cause by this program. Please use responsibly.

The material contained in this repository is licensed under GNU GPLv3.
