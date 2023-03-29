package specs

type Payload interface {
	Table() string
	Database() string
	Index() int

	Fields() []DriverField
	Join() []DriverJoin
	Where() []DriverWhere

	SetFields([]DriverField) Payload
	SetJoins([]DriverJoin) Payload
	SetWheres([]DriverWhere) Payload

	Mapping() ([]any, error)
	OnScan([]any) error
}

type PayloadAugmented[T Model] interface {
	Payload
	Result() []T
}
