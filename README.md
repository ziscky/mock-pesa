# mock-pesa
[![Docker Repository on Quay](https://quay.io/repository/ziscky/mock-pesa/status)](https://quay.io/repository/ziscky/mock-pesa)
[![Build Status](https://goreportcard.com/badge/github.com/ziscky/zist)](https://goreportcard.com/report/github.com/ziscky/mock-pesa)
[![Build Status](https://travis-ci.org/ziscky/zist.svg?branch=master)](https://travis-ci.org/ziscky/mock-pesa)
[![Issue Count](https://codeclimate.com/github/ziscky/mock-pesa/badges/issue_count.svg)](https://codeclimate.com/github/ziscky/mock-pesa)


A set of mock MPESA APIs designed to work exactly like the real APIs. Great for unit tests and payment system hardening

### Precompiled binaries

Precompiled binaries for released versions are available in the
[*releases* section](https://github.com/ziscky/mock-pesa/releases)
of the GitHub repository. Supported OS/Arch:

 1. Darwin X64
 2. FreeBSD X64
 3. Linux X64
 4. Windows X64

 

### Container images

Container images are available on https://quay.io/repository/ziscky/mock-pesa.  
Simply: `docker pull quay.io/ziscky/mock-pesa`  

### Getting Started
For the docker container:  
`docker run -e MERCHANT_ID=1234 -e PASSKEY=4321 -p 7000:7000 --rm quay.io/ziscky/mock-pesa `  

For the precompiled binaries:  
`MERCHANT_ID=1234 PASSKEY=4321 ./mock-pesa`  
OR  
`./mock-pesa -conf=/path/to/conf`  

### Example Config
`
MaxAmount=70000  

MinAmount=10  

MaxCustomerTransactionPerDay=150000  

MerchantID="12345"  

CallBackDelay=0  

SAGPasskey="54321"  

EnabledAPIS = ["c2b"]  
`

### Building From Source
`go get github.com/ziscky/mock-pesa`  
`go test`  
`go build github.com/ziscky/mock-pesa -o mock-pesa`   

### How To
Follow the official MPESA API guide, replace:  
`Endpoint: http://localhost:7000/`  
`WSDL: http://localhost:7000/wsdl/get`  
That's right, works exactly the same  

For custom scenarios e.g `Customer with insufficient funds`:  
`Endpoint: http://localhost:7000/{code}`  
where code is one of the official MPESA response codes,in this case: `01`  

### Contiributing
I'm very open to PRs.  

 - Fork
 - Create Branch
 - Do magic
 - Initiate PR

