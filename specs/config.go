package specs

type Config interface {
	Name() string
	SetName(name string) Config

	Driver() string
	SetDriver(driver string) Config

	User() string
	SetUser(user string) Config

	Password() string
	SetPassword(password string) Config

	Host() string
	SetHost(host string) Config

	Port() int
	SetPort(port int) Config

	Database() string
	SetDatabase(database string) Config

	Locale() string
	SetLocale(locale string) Config
}
