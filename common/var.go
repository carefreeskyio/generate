package common

var ServiceFileName = "service.go"

var RpcFileName = "rpc.go"

var ServiceTemp = `package service

{import_list}

{service_struct}

func NewService() *Service {
	{new_service}
}

{service_func}
`

var ServiceFuncTemp = `func (s *Service) {name} (ctx context.Context, {param}) (err error) {
	return s.{service_name}.{name}(ctx, {pass_through})
}

`

var RpcTemp = `package rpc

{import_list}

type Client struct {
	XClient *rpcxclient.Client
}

var (
	c    *Client
	once sync.Once
)

func NewClient() (*Client, error) {
	var (
		err        error
		rpcXClient *rpcxclient.Client
	)

	if c == nil {
		once.Do(func() {
			rpcXClient, err = rpcxclient.NewClient(initOptions())
			if err != nil {
				return
			}
			c = &Client{
				XClient: rpcXClient,
			}
		})
	}

	return c, err
}

// 获取初始化rpcXClient客户端属性，可根据实际需求修改
func initOptions() (options rpcxclient.Options) {
	options = rpcxclient.DefaultOptions
	options.BasePath = "{base_path}"
	options.ServerName = "{service_name}"
	options.Addr = strings.Split(global.Config.Registry.Addr, " ")
	options.Group = global.Config.Registry.Group
	options.Timeout = time.Duration(global.Config.Rpc.WithTimeout) * time.Second

	return options
}

{rpc_func}
`

var RpcFuncTemp = `func (c *Client) {name}(ctx context.Context, {param}) (err error) {
	return c.XClient.Call(ctx, "{name}", {pass_through})
}

`
