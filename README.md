# puredns
![puredns](https://user-images.githubusercontent.com/55468528/88924057-b420a680-d240-11ea-8ff7-6fb988db8ec1.png)

[massdns](https://github.com/blechschmidt/massdns) is an incredibly powerful DNS resolver used to perform bulk lookups. With the proper bandwidth and a good list of public resolvers, it can resolve millions of queries in just a few minutes.

Unfortunately, the results are only as good as the answers provided by the public resolvers used. They are often polluted by DNS poisoned entries. Wildcard subdomains are also a pain to deal with, as they add a lot of noise to the list of resolved subdomains.

**puredns** is a bash and python application that uses massdns to _accurately_ perform DNS bruteforcing and mass resolving. It ensures the results obtained by public resolvers are clean and can work around DNS poisoning by validating the answers obtained using a list of trusted resolvers. It also handles wildcard subdomains correctly through its companion Python script, **wildcarder**.

Think this is useful? :star: Star us on GitHub — it helps!

## Features

* Accurately resolves thousands of DNS queries per second using massdns and a list of public DNS resolvers
* Supports DNS bruteforcing using a wordlist and a root domain
* Cleans wildcard root subdomains with a minimal number of queries
* Validates that the results are free of DNS poisoning by running against a list of known, trustable resolvers
* Saves list of valid subdomains, wildcard root domains and answers, and clean massdns output containing only the valid entries
* Comes with a massive list of tested public resolvers from [public-dns.info](https://public-dns.info/)

## How it works

![puredns](https://user-images.githubusercontent.com/55468528/88879682-ee665580-d1f8-11ea-9239-eb895790aa63.gif)

In the image above, you can see puredns in action against the domain store.yahoo.com using a small wordlist of the 100k most common subdomains. This is happening in real time.

As part of its workflow, puredns perfoms four steps automatically:

### 1. Preparation and sanitization of domains to resolve

When in resolve mode, puredns expects a list of domains to resolve.

When in bruteforce mode, puredns creates the list of domains to resolve from the domain provided and a wordlist.

In both cases, the list of domains is sanitized. Only entries containing valid characters that can be found in a domain name are kept (essentially `[a-zA-Z0-9\.\-]`).

This step can be skipped using the `--skip-sanitize` flag.

### 2. Dumb mass resolve

puredns will then perform a mass resolve of all the domains in the list using massdns. It saves the results to a temporary file.

The results are usually polluted: some public resolvers will have sent poisoned replies, and wildcard subdomains can also quickly inflate the results.

This step is mandatory and cannot be skipped.

### 3. Wildcard detection

The application then uses its companion _wildcarder_ Python script to detect and extract all the wildcard root subdomains from the massdns results file.

wildcarder will reuse the massdns output to minimize the number of queries it needs to perform. It is quite common for wildcarder to figure out the structure of wildcard subdomains with as few as 5-10 DNS queries when it has a massdns file to prime its cache.

puredns will then clean up all the wildcard subdomains from its results and keep only the wildcard root subdomains that resolve correctly.

This step can be skipped using the `--skip-wildcard-check` flag.

### 4. Validation

To protect against DNS poisoning, puredns uses massdns one last time to validate the remaining results using a list of trusted DNS resolvers.

This step is done at a slower pace as not to hit any rate limiting on the trusted resolvers. We try to limit the rate to 10 queries per second per resolver.

This step can be skipped using the `--skip-validation` flag.

Once this is done, the resulting files are clean of wildcard subdomains and DNS poisoned answers.

Only the resulting valid domains are sent to stdout so that you can pipe the results to other programs. The rest of the information output by puredns is sent to stderr.

## Installation

puredns requires massdns to be installed on the host machine. [Follow the instructions](https://github.com/blechschmidt/massdns#compilation) to compile massdns on your system.

If the path to the massdns binary is present in the PATH environment variable, puredns will work out of the box. On most systems, a good place to copy the massdns executable is `/usr/local/bin`.

Otherwise, you will need to specify the path to the massdns binary file using the `--bin` command line argument.

The script also requires a few other dependencies:

* python3
* python3-dnspython
* pv

To ensure the dependencies are installed on Ubuntu, you can use the following command line:

```
sudo apt install -y python3 python3-dnspython pv
```



Once the dependencies are installed, simply clone the repository to start using puredns:

```
git clone https://github.com/d3mondev/puredns.git
```

## Usage

### puredns

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
                -b,  --bin <path>                       Path to massdns binary file

                -r,  --resolvers <filename>             Text file containing resolvers
                -rt, --resolvers-trusted <filename>     Text file containing trusted resolvers

                -l,  --limit                            Limit queries per second for public resolvers
                                                        (default: unlimited)
                -lt, --limit-trusted                    Limit queries per second for trusted resolvers
                                                        (default: 10 * number of trusted resolvers)

                -ss, --skip-sanitize                    Do not sanitize the list of domains to test
                                                        By default, domains are set to lowercase and
                                                        only valid characters are kept
                -sw, --skip-wildcard-check              Do no perform wildcard detection and filtering
                -sv, --skip-validation                  Do not validate massdns results using trusted resolvers

                -w,  --write <filename>                 Write valid domains to a file
                -wm, --write-massdns <filename>         Write massdns results to a file
                -ww, --write-wildcards <filename>       Write wildcard root subdomains to a file
                -wa, --write-answers <filename>         Write wildcard DNS answers to a file

                -h, --help                              Display this message
```

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

## Examples

![puredns_yahoo](https://user-images.githubusercontent.com/55468528/88879692-f1f9dc80-d1f8-11ea-9e96-70a9ef107246.png)

### Subdomain bruteforcing

Using the included list of resolvers, puredns can bruteforce a massive list of subdomains using a wordlist named `all.txt`:

`puredns bruteforce all.txt domain.com`

### Resolving a list of domains

You can also resolve a list of domains contained in a text file (one per line).

`puredns resolve domains.txt`

### Saving the results to files

You can save the following information to files to reuse it in your workflows:

* **domains**: clean list of domains that resolve correctly
* **wildcard root domains**: List of the wildcard root domains found (ie. *\*.store.yahoo.com*)
* **wildcard answers**: list of DNS answers given by the wildcard subdomains (A and CNAME records)
* **massdns results file (simple text output)**: can be used as a reference and to extract A and CNAME records.

```
puredns resolve domains.txt --write valid_domains.txt \
                            --write-wildcard-answers wildcard_answers.txt \
                            --write-wildcards wildcards.txt \
                            --write-massdns massdns.txt
```

### Using custom resolvers

You can use a custom list of resolvers with puredns. Simply pass the `-r` argument to the script.

You can also specify a list of custom trusted resolvers with the `-rt` argument.

`puredns resolve domains.txt -r resolvers.txt -rt trusted.txt`

## Resources

[shuffleDNS](https://github.com/projectdiscovery/shuffledns) is a good alternative written in go that handles wildcard subdomains using a different algorithm.

[public-dns.info](https://public-dns.info/) continuously updates a list of public and free DNS resolvers.

[DNS Validator](https://github.com/vortexau/dnsvalidator) can be used to curate your own list of public DNS resolvers.

[all.txt wordlist](https://gist.github.com/jhaddix/f64c97d0863a78454e44c2f7119c2a6a) Jhaddix's iconic `all.txt` wordlist commonly used for subdomain enumeration.

## Author

I'm d3mondev, a coder and security enthusiast. [Follow me on twitter](https://twitter.com/d3mondev)!

## Disclaimer & License

The resolvers included in this repository are present for reference only. The author is not responsible for any misuse of the resolvers in that list. It is the user's responsibility to curate a list of resolvers you are authorized to use.

Usage of this program for attacking targets without consent is illegal. It is the user's responsibility to obey all applicable laws. The developer assumes no liability and is not responsible for any misuse or damage cause by this program. Please use responsibly.

The material contained in this repository is licensed under GNU GPLv3.
