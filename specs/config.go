package specs

type Config interface {
	Name() string
	SetName(name string)

	Driver() string
	SetDriver(driver string) Config

	User() string
	SetUser(user string)

	Password() string
	SetPassword(password string)

	Host() string
	SetHost(host string)

	Port() int
	SetPort(port int)

	Database() string
	SetDatabase(database string)

	Locale() string
	SetLocale(locale string)
}
