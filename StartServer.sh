# ./STOIC run server --ip=127.0.0.1 --port=5001 --batch=24 | tee -a edge_log.txt
# ./STOIC run server --ip=127.0.0.1 --port=5001 --image=132 --batch=1 | tee -a cpu_log.txt
# ./STOIC run server --ip=127.0.0.1 --port=5001 --batch=24 | tee -a gpu1_log.txt
# ./STOIC run server --ip=127.0.0.1 --port=5001 --batch=3 --preset=true | tee -a gpu2_log.txt
./STOIC run server --ip=127.0.0.1 --port=5001 --image=1 | tee -a STOIC_log.txt