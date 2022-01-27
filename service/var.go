package service

var fileName = "service.go"

var fileTemp = `package main

{import_list}

{service_struct}

func NewService() *Service {
	{new_service}
}

{service_func}
`

var funcTemp = `func (s *Service) {name} (ctx context.Context, {param}) (err error) {
	return s.{service_name}.{name}(ctx, {pass_through})
}

`
