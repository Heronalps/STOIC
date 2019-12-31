# ./stoic run server --ip=127.0.0.1 --port=5001 --batch=24 | tee -a edge_log.txt
# ./stoic run server --ip=127.0.0.1 --port=5001 --image=132 --batch=1 | tee -a cpu_log.txt
./stoic run server --ip=127.0.0.1 --port=5001 --batch=24 | tee -a gpu1_log.txt
# ./stoic run server --ip=127.0.0.1 --port=5001 --batch=3 --preset=true | tee -a gpu2_log.txt
# ./stoic run server --ip=127.0.0.1 --port=5001 --batch=24 | tee -a stoic_log.txt