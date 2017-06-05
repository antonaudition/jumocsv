# Jumo CSV Collator

Every two months we are given a text file from our accounting department detailing all the loans
given out on our networks. This needs to be a validated against our internal systems by the tuple
of (Network, Product, Month). We will drill down as necessary into problem areas.

The expectation is a file called Output.csv with a line detailing the totals by the tuple of
(Network, Product, Month) summing the amount and count of loans.

## Getting Started

After cloning the project the easiest way to get started is by running the make commands in the 
root of the project

### Prerequisites

The following tools are required to build and run this project

```
go 1.8 (https://golang.org/)
```

### Installing

The source compiles to a single binary `jumocsv`. If you are familiar with the go build tools, you can 
clone the repo to your `$GOPATH` and build it as normal. Otherwise the simplest way to get the binary is
 to run the following commands after installing go.

```
go get github.com/antonaudition/jumocsv
go install github.com/antonaudition/jumocsv/cmd/jumocsv
```

### Usage

The binary takes the input csv as an argument and produces a file `Output.csv`, ex.

```
$ jumocsv Loans.csv
$ ls
> Loans.csv Output.csv
```
