## pathcrawler

Discover new paths via scanning html.

## Installation
```
go install github.com/rix4uni/pathcrawler@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/pathcrawler/releases/download/v0.0.1/pathcrawler-linux-amd64-0.0.1.tgz
tar -xvzf pathcrawler-linux-amd64-0.0.1.tgz
rm -rf pathcrawler-linux-amd64-0.0.1.tgz
mv pathcrawler ~/go/bin/pathcrawler
```
Or download [binary release](https://github.com/rix4uni/pathcrawler/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/pathcrawler.git
cd pathcrawler; go install
```

## Usage
```
Usage of pathcrawler:
  -complete-url
        Complete relative URLs based on the input URL
  -concurrent int
        Number of concurrent requests (default 50)
  -delay duration
        Delay between requests
  -only-complete
        Show only complete URLs starting with http:// or https://
  -silent
        silent mode.
  -timeout duration
        Timeout duration for HTTP requests (default 10s)
  -version
        Print the version of the tool and exit.
```

## Examples
Single Target:
```
▶ echo "https://dell.com" | pathcrawler -silent
```

Multiple Targets:
```
▶ cat targets.txt
https://google.com
https://dell.com

▶ cat targets.txt | pathcrawler -silent
```