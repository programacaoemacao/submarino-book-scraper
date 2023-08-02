package interfaces

type Exporter interface {
	Export(items interface{}) error
}
