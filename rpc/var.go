package rpc

var filePath = "./rpc"

var importList = map[string]byte{
	"\t\"context\"":                            0,
	"\t\"github.com/carefreex-io/config\"":     0,
	"\t\"github.com/carefreex-io/rpcxclient\"": 0,
	"\t\"sync\"":                               0,
	"\t\"time\"":                               0,
}

var fileTemp = `package rpc

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
			rpcXClient, err = rpcxclient.NewClient(getOptions())
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
func getOptions() *rpcxclient.Options {
	options := rpcxclient.DefaultOptions
	options.RegistryOption.BasePath = "{base_path}"
	options.RegistryOption.ServerName = "{service_name}"
	options.RegistryOption.Addr = config.GetStringSlice("Registry.Addr")
	options.RegistryOption.Group = config.GetString("Registry.Group")
	options.Timeout = config.GetDuration("Rpc.Timeout") * time.Second

	return options
}

{rpc_func}
`

var funcTemp = `func (c *Client) {name}(ctx context.Context, {param}) (err error) {
	return c.XClient.Call(ctx, "{name}", {pass_through})
}

`
