module go_script/grpc_2/c

go 1.16

replace go_script/grpc_2/pb => ../pb

require (
	go_script/grpc_2/pb v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	google.golang.org/grpc v1.46.0
)
