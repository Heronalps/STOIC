SELECT 'epoch', 'image_num', 'runtime', 'pred_total', 'pred_transfer', 'pred_deploy', 'pred_proc', 'act_total', 'act_transfer', 'act_deploy', 'act_proc'
UNION
SELECT UNIX_TIMESTAMP(time_stamp) as epoch, image_num, runtime, pred_total, pred_transfer, pred_deploy, pred_proc,
act_total, act_transfer, act_deploy, act_proc from LogTime where app="image-clf-inf" and version="1.1"
INTO OUTFILE '/var/lib/mysql-files/logTime.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n';