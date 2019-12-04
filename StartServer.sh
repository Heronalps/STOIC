go run main.go -node=server -ip=128.111.45.118 -port=5001 -image=20 -runtime=edge -batch=10 | tee -a edge_log.txt
go run main.go -node=server -ip=128.111.45.118 -port=5001 -image=20 -runtime=cpu -batch=10 | tee -a cpu_log.txt
go run main.go -node=server -ip=128.111.45.118 -port=5001 -image=20 -runtime=gpu1 -batch=10 | tee -a gpu1_log.txt
go run main.go -node=server -ip=128.111.45.118 -port=5001 -image=20 -runtime=gpu2 -batch=10 | tee -a gpu2_log.txt
go run main.go -node=server -ip=128.111.45.118 -port=5001 -image=20 -batch=10 | tee -a stoic_log.txt