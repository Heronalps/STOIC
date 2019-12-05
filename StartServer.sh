# go run main.go -node=server -ip=127.0.0.1 -port=5001 -runtime=edge -batch=24 | tee -a edge_log.txt
go run main.go -node=server -ip=127.0.0.1 -port=5001 -runtime=cpu -batch=24 | tee -a cpu_log.txt
# go run main.go -node=server -ip=127.0.0.1 -port=5001 -runtime=gpu1 -batch=24 | tee -a gpu1_log.txt
# go run main.go -node=server -ip=127.0.0.1 -port=5001 -runtime=gpu2 -batch=24 | tee -a gpu2_log.txt
# go run main.go -node=server -ip=127.0.0.1 -port=5001 -batch=24 | tee -a stoic_log.txt