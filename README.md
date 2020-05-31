# COVID 19 API

The COVID 19 API is organized around REST. Provided API has predictable resource-oriented URLs, accepts form-encoded request bodies, returns JSON-encoded responses, and uses standard HTTP response codes, verbs. The api is driven by the **Covid 19 dataset** from [**Novel Coronavirus (COVID-19) Cases, provided by JHU CSSE**](https://github.com/CSSEGISandData/COVID-19)

<!-- Available here:  -->

***
## Table of contents
- [Getting Started](#getting-started)
    - [Golang Installation on OS X](#golang-installation-on-os-x)
    - [Golang Installation on Linux](#golang-installation-on-linux)
    - [Golang Installation on Windows](#golang-installation-on-windows)
    - [Postman Installation](#postman-installation)
- [Installation](#installation)
- [Download Dataset](#download-dataset)
- [Build and Run](#build-and-run)
    - [Localhost](#localhost)
    - [Docker](#docker)
- [Go Packages](#go-packages)
- [Resources](#resources)
- [Areas of improvements](#areas-of-improvements)
- [Contributing](#contributing)
    - [Known Issues](#known-issues)
- [Things I hate about Golang](#things-i-hate-about-golang)
- [Credits](#credits)
- [License](#license)

## Getting Started
For development, you will only need [Golang](https://golang.org/) installed on your environement. 
**Note:** This project is developed on *Golang version 1.13*

#### Golang installation on OS X

You will need to use a Terminal. On OS X, you can find the default terminal in
`/Applications/Utilities/Terminal.app`.

Please install [Homebrew](http://brew.sh/) if it's not already done with the following command.

    $ ruby -e "$(curl -fsSL https://raw.github.com/Homebrew/homebrew/go/install)"

If everything when fine, you should run

    brew update
    brew install golang
    
    - Then add those lines to export the required variables
    export GOPATH=$HOME/go-workspace # don't forget to change your path correctly!
    export GOROOT=/usr/local/opt/go/libexec

    - Verifying the Go Installation
    go version

#### Golang installation on Linux

    wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
    sha256sum go1.13.linux-amd64.tar.gz # verify the Go tarball
    sudo tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz

    - Adjusting the Path Variable (appending the following line to the /etc/profile file (for a system-wide installation) or the $HOME/.profile file (for a current user installation))
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go  # don't forget to change your path correctly!

    - Save the file, and load the new PATH environment variable into the current shell session
    source ~/.profile

    - Verifying the Go Installation
    go version

#### Golang installation on Windows

Just go on [official Golang website](https://golang.org/doc/install?download=go1.14.2.windows-amd64.msi) & grab the installer.

#### Postman Installation

Download from [here](https://www.postman.com/downloads/) for your platform and install it.

## Installation
```
    $ git clone https://github.com/ParthoShuvo/covid-19-api.git
    $ cd covid-19-api
    $ export GOFLAGS=-mod=vendor
    $ export GO111MODULE=on
    $ go mod download
    $ go mod vendor
    $ go mod verify
```

## Download [Dataset](https://github.com/CSSEGISandData/COVID-19)
**open bash/zsh shell** 
```
    $ chmod +x covid-19-dataset.sh
    $ ./covid-19-dataset.sh
```

## Build And Run
**Note:** Please change server *bind-address* and *port* at **covid_19_api.json** file before run following.

#### Localhost
```
    $ go build
    $ ./covid-19-api
```

#### Docker
```
    $ docker build -f dev.Dockerfile -t covid-19-api .
    $ docker run -it --rm -p ${host_port}:${docker_port} covid-19-api
```
## Go Packages

- **Logging** [logrus](https://github.com/sirupsen/logrus)
- **Http Router and URL Matcher**  - [gorilla/mux](https://github.com/gorilla/mux)
- **CSV Serializer & Deserializer** - [gocsv](https://github.com/gocarina/gocsv)
- **Functional Programming in Go** - [fpingo](https://github.com/ParthoShuvo/fpingo) :construction: (Hope to get some motivation to complete :blush:)

## Resources
- [**Go In Action**](https://www.manning.com/books/go-in-action)
- [**Practical Persistence in Go: Organising Database Access**](https://www.alexedwards.net/blog/organising-database-access)
- [**Debugging Go Code using VSCode**](https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code)
- [**Golang Error Handling — Best Practice**](https://itnext.io/golang-error-handling-best-practice-a36f47b0b94c)
- [**Awesome Go**](https://awesome-go.com/)

## Areas of improvements
*Need to work on following to make code more idiomatic and clean*
- Apply [**Effective Go Rules**](https://golang.org/doc/effective_go.html)
- Efficient logging & error handling
- Improve dataset reading and searching performance 
- Rest API endpoints design improvement
- Follow [**Go in Practice**](https://www.manning.com/books/go-in-practice)


## Contributing
Pull requests are welcome. Please make sure that your PR is [well-scoped](https://www.netlify.com/blog/2020/03/31/how-to-scope-down-prs/).
For major changes, please open an issue first to discuss what you would like to change. 

## Things I hate about Golang

*Following links tell you why I hate*
- [**Go is not good**](https://github.com/ksimka/go-is-not-good)
- [**Go good bad ugly**](https://bluxte.net/musings/2018/04/10/go-good-bad-ugly/)
- [**50 Shades of Go**](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)

### Known issues
Visit [issues](https://github.com/ParthoShuvo/covid-19-api/issues) section.

### Contributors
<table>
  <tr>
    <td align="center"><a href="https://www.linkedin.com/in/parthoshuvo/"><img src="https://avatars3.githubusercontent.com/u/9255705?s=460&u=15a0c89028fcfe11868a679406e90ef94eeeedd9&v=4" width="200px;" alt=""/><br /><sub><b>Shuvojit Saha</b></sub></a><br /><a href="https://github.com/ParthoShuvo/covid-19-api/commits?author=ParthoShuvo" title="Code">💻</a> <a href="#infra-sruti" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/ParthoShuvo/covid-19-api/issues/created_by/ParthoShuvo" title="Bug reports">🐛</a><a href="#ideas-sruti" title="Ideas, Planning, & Feedback">💡</a></td>
    </tr>
</table>

### Credits
- Kudos to [**Steef de Rooi**](https://www.linkedin.com/in/steefderooi/?originalSubdomain=nl) for his idiomatic golang code. And so I started learning Golang and being always inspired from his clean code.
- [**Global COVID-19 Dataset**](https://github.com/CSSEGISandData/COVID-19) by **Johns Hopkins University**
- **Readme.md** format inspired from [sruti](https://github.com/sruti/covid19-riskfactors-app)
- Commit message emoji from [**Gitmoji**](https://gitmoji.carloscuesta.me/)


## License
[MIT](https://choosealicense.com/licenses/mit/)


