# msbingo
[![Build Status](https://travis-ci.org/Nordstrom/msbingo.svg?branch=master)](https://travis-ci.org/Nordstrom/msbingo)
[![Coverage Status](https://coveralls.io/repos/github/Nordstrom/msbingo/badge.svg?branch=master)](https://coveralls.io/github/Nordstrom/msbingo?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nordstrom/msbingo)](https://goreportcard.com/report/github.com/Nordstrom/msbingo)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![](https://img.shields.io/badge/windows-ready-green.svg)

| NBFX | NBFS |
|:----:|:----:|
|[![GoDoc](https://godoc.org/github.com/Nordstrom/msbingo/nbfx?status.png)](https://godoc.org/github.com/Nordstrom/msbingo/nbfx)|[![GoDoc](https://godoc.org/github.com/Nordstrom/msbingo/nbfs?status.png)](https://godoc.org/github.com/Nordstrom/msbingo/nbfs)

Need to build a golang service that can interop with a .NET service which only speaks WCF Binary Soap? Msbingo is that missing piece of the puzzle.

Msbingo is an implementation of NBFX and NBFS, which is expressed in HTTP terms as `Content-Type: application/soap+msbin1`, written in pure Go to enable interop with a WCF service and Go with no other dependencies (like Mono or Windows).

This implementation follows the Microsoft specification closely in naming, structure and example. Tests for both decoding and encoding have been written for each of the Structure Examples provided, to validate the given bytes encode to the given XML and vice versa.

Currently the decoding side is more (nearly!) complete in terms of implementation of the individual record types. You can see the current state by running the tests.

Encoding is less complete, as it is not necessary to encode exactly as .NET WCF would. It is only necessary to encode properly, such that the target service can decode the original XML message properly.

Feel free to fill out the codec on either the encoding or decoding sides by fixing additional tests or writing more.

# Usage

``` go
url := fmt.Sprintf("%s/Path/To/ExampleService.svc", s.apiBaseUrl)

xmlInput := "<s:Envelope xmlns:a=\"http://www.w3.org/2005/08/addressing\" xmlns:s=\"http://www.w3.org/2003/05/soap-envelope\"><s:Header><a:Action s:mustUnderstand=\"1\">action</a:Action></s:Header><s:Body><Inventory>0</Inventory></s:Body></s:Envelope>"

encodedXml, err := nbfs.NewEncoder().Encode(bytes.NewBufferString(xmlInput))
if err != nil {
	// handle encoding error
}

req, err := http.NewRequest("POST", url, bytes.NewBuffer(encodedXml))
if err != nil {
	// handle request creation error
}

req.Header.Add("Content-Type", "application/soap+msbin1")
resp, err := httpClient.Do(req)
if err != nil {
	// handle response error
}
defer resp.Body.Close()

xmlRes, err := nbfs.NewDecoder().Decode(resp.Body)
if err != nil {
	// handle decoding error
}

// do something with your decoded xml response
```

# Background
Application/soap+msbin1 encoding was a blocking issue for modernizing services from WCF to platform-agnostic technologies such as Go. We needed to be able to make calls to dependency services that spoke msbin1 and were not going to be updated or even reconfigured, but we did not want to introduce unnecessary complexity such as workarounds like .NET-based WCF request translator proxies or deploying Mono with our service instances. Initially we tried the Mono deployment route, which, while it would have worked well enough, significantly complicated our deployment pipeline, thus erasing one of the major advantages of golang.

This project began as an unstructured/personal time project for a few of our team members to help accelarate our golang expertise and to "geek out" a bit on some lower level code than we normally work on. We decided implementing our own codec in pure golang from the MS-published spec was worth the time and we were right. We are open sourcing this in the hope of helping other groups out there in similar circumstances.

# Codec details
NBFS is a codec developed by Microsoft for use primarily by WCF webservices. It is essentially a binary encoding for Soap XML messages optimized for reducing bytes over the wire.  They have published the specification in multiple parts:
* [NBFX (.NET Binary Format: XML Data Structure)](https://msdn.microsoft.com/en-us/library/cc219210.aspx)
* [NBFS (.NET Binary Format: SOAP Data Structure)](https://msdn.microsoft.com/en-us/library/cc219175.aspx)
where NBFS is essentially NBFX with standard DictionaryString entries for strings commonly used in SOAP, such as "Envelope", "http://www.w3.org/2003/05/soap-envelope/", etc., to minimize the bytewise size overhead of the SOAP protocol.

# Contributors

* [Khoa Nguyen (khoad)](https://github.com/khoad/)
* [Tim Erickson (timerickson)](https://github.com/timerickson/)
* [Ryan Brewer (ryanthecubfan)](https://github.com/ryanthecubfan/)

# License

msbingo is released under the Apache-2.0 license.
