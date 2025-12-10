module zyg/datasource

go 1.14

require (
	zyg/datamodels v0.0.0 // indirect
	zyg/repositories v0.0.0 // indirect
	zyg/service v0.0.0 // indirect
	zyg/web v0.0.0 // indirect
)

replace (
	zyg/datamodels => ../datamodels
	zyg/datasource => ../datasource
	zyg/repositories => ../repositories
	zyg/service => ../service
	zyg/web => ../web
)
