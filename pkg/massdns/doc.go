// Package massdns provides a Resolver object used to invoke the massdns binary file.
//
// The package contains a LineReader struct that implements the io.Reader interface. It is used
// to read strings line by line from an io.Reader while throttling the results according to a
// rate-limit specified. The LineReader is passed to the stdin of massdns, allowing it to
// approximately respect the number of DNS queries per second wanted.
package massdns
