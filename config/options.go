package config

//Option optional parameter structure
type Option struct {
	name  string
	value interface{}
}

//NewOption new Option
func NewOption(name string, value interface{}) *Option {
	return &Option{
		name:  name,
		value: value,
	}
}

//Name of the Option
func (o *Option) Name() string {
	return o.name
}

//Value of the Option
func (o *Option) Value() interface{} {
	return o.value
}
