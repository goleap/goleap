package drivers

import (
	"database/sql"
	"reflect"
	"sync"
)

type DbCacheConnection interface {
	RegisterConn(cnx *sql.Conn, connectionId int64)
	GetConnectionId(cnx *sql.Conn) int64
}

type dbCacheConnection struct {
	sync.Mutex
	connections map[uintptr]int64
}

var DbCacheConnectionInstance DbCacheConnection

func init() {
	DbCacheConnectionInstance = &dbCacheConnection{
		connections: make(map[uintptr]int64),
	}
}

func (d *dbCacheConnection) extractPointerFromCnx(cnx *sql.Conn) uintptr {
	rf := reflect.ValueOf(cnx)
	return rf.Elem().FieldByName("dc").Pointer()
}

func (d *dbCacheConnection) RegisterConn(cnx *sql.Conn, connectionId int64) {
	d.Lock()
	defer d.Unlock()

	ptr := d.extractPointerFromCnx(cnx)
	d.connections[ptr] = connectionId
}

func (d *dbCacheConnection) GetConnectionId(cnx *sql.Conn) int64 {
	d.Lock()
	defer d.Unlock()

	ptr := d.extractPointerFromCnx(cnx)

	return d.connections[ptr]
}
