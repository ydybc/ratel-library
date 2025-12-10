module zyg/web

go 1.14

require (
	github.com/BelieveR44/goseaweedfs v0.0.0-20200531112332-37d0726cea2d
	github.com/BurntSushi/toml v0.3.1
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394 // indirect
	github.com/dchest/captcha v0.0.0-20170622155422-6a29415a8364
	github.com/go-redis/redis v6.15.8+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/securecookie v1.1.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/kataras/iris/v12 v12.1.8
	github.com/kirinlabs/HttpRequest v1.0.5 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/urfave/cli v1.22.4
	gopkg.in/olivere/elastic.v5 v5.0.85
	zyg/datamodels v0.0.0
	zyg/datasource v0.0.0 // indirect
	zyg/repositories v0.0.0
	zyg/service v0.0.0
)

replace (
	github.com/linxGnu/goseaweedfs => ../../github.com/linxGnu/goseaweedfs
	zyg/datamodels => ../datamodels
	zyg/datasource => ../datasource
	zyg/repositories => ../repositories
	zyg/service => ../service
	zyg/web => ../web
)
