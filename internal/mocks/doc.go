/*
Package mocks will have all the mocks of the library.
*/
package mocks // import "github.com/kiritoxkiriko/go-http-metrics/internal/mocks"

//go:generate mockery -output ./metrics -outpkg metrics -dir ../../metrics -name Recorder
//go:generate mockery -output ./middleware -outpkg middleware -dir ../../middleware -name Reporter
