// Package procreader provides a ProcReader object that implements the io.Reader interface and generates
// its data from a user-specified callback function.
//
// Use the New function to create a new ProcReader and pass it a callback function that is invoked to
// create new data when the ProcReader's buffers are empty.
package procreader
