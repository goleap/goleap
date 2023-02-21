package specs

type Payload interface {
	Table() string
	Database() string
	Index() int
	Fields() []DriverField
	Join() []DriverJoin
	Where() []DriverWhere
	Mapping() []any
	OnScan([]any) error
}
