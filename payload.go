package goleap

import (
	"github.com/goleap/goleap/connector/driver"
	"github.com/goleap/goleap/schema"
)

type payload[T schema.Model] struct {
	*orm[T]

	focusedSchemaFields      []schema.Field
	focusedDriverFields      []driver.Field
	requiredJoins            []driver.Join
	focusedSchemaValueFields []any

	result []T
}

func (p *payload[T]) Database() string {
	return p.orm.schema.DatabaseName()
}

func (p *payload[T]) Index() int {
	return p.orm.schema.Index()
}

func (p *payload[T]) Fields() []driver.Field {
	return p.orm.getFocusedDriverFields()
}

func (p *payload[T]) Join() []driver.Join {
	return p.orm.getRequiredJoins()
}

func (p *payload[T]) Where() []driver.Where {
	return []driver.Where{}
}

func (p *payload[T]) ResultType() []any {
	return p.orm.getFieldsTypeSchemaToDriver()
}

func (p *payload[T]) OnScan(result []any) (err error) {
	sch := p.schema.Copy()
	for i, field := range p.orm.focusedSchemaFields {
		sch.GetFieldByName(field.RecursiveFullName()).Set(result[i])
	}
	p.result = append(p.result, sch.Get().(T))
	return
}

func (p *payload[T]) Table() string {
	return p.orm.schema.TableName()
}

func newPayload[T schema.Model](orm *orm[T]) driver.Payload {
	p := new(payload[T])
	p.orm = orm
	return p
}
