module go_script/grpc_1/s

go 1.16

replace github.com/mangenotwork/man/core/grpc_1/pb => ../pb

require (
	github.com/mangenotwork/man/core/grpc_1/pb v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.46.0
)
