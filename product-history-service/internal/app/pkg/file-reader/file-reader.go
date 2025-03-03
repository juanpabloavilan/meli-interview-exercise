package filereader

import (
	"encoding/csv"

	gocsv "github.com/gocarina/gocsv"
)

type CSVRowUnmarshallerErr struct {
	Err error
}

func (e CSVRowUnmarshallerErr) Error() string {
	return e.Err.Error()
}

type CSVRowUnmarshaller[T any] interface {
	ReadUnmarshalCSVRow() (any, error)
}

type csvRowUnmarshaller struct {
	Unmarshaller *gocsv.Unmarshaller
}

func NewCSVRowUnmarshaller[T any](r *csv.Reader) (CSVRowUnmarshaller[T], error) {
	u, err := gocsv.NewUnmarshaller(r, new(T))
	if err != nil {
		return nil, err
	}

	return &csvRowUnmarshaller{
		Unmarshaller: u,
	}, nil
}

func (r csvRowUnmarshaller) ReadUnmarshalCSVRow() (any, error) {
	out, err := r.Unmarshaller.Read()
	if err != nil {
		return nil, err
	}

	return out, nil
}
