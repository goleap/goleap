package specs

type NewPayload[T Model] func(model ...Model) PayloadAugmented[T]

type Payload interface {
	Table() string
	Database() string
	Index() int

	Fields() []DriverField
	Join() []DriverJoin
	Where() []DriverWhere
	Limit() DriverLimit

	SetFields([]DriverField) Payload
	SetJoins([]DriverJoin) Payload
	SetWheres([]DriverWhere) Payload
	SetLimit(DriverLimit) Payload

	Mapping() ([]any, error)
	OnScan([]any) error
}

type PayloadAugmented[T Model] interface {
	Payload
	Result() []T
}
