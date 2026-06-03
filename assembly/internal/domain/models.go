package domain

type Car struct {
	VIN   string
	Brand string
	Year  int32
}
type Engine struct {
	ID         string
	Horsepower int32
}
type Transmission struct {
	ID   string
	Type string
}
type CarSpec struct {
	Car          Car
	Engine       Engine
	Transmission Transmission
}
