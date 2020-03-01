SELECT 'epoch', 'image_num', 'runtime', 'pred_total', 'pred_transfer', 'pred_deploy', 'pred_proc', 'act_total', 'act_transfer', 'act_deploy', 'act_proc'
UNION
SELECT UNIX_TIMESTAMP(time_stamp) as epoch, image_num, runtime, pred_total, pred_transfer, pred_deploy, pred_proc,
act_total, act_transfer, act_deploy, act_proc from LogTime where app="image-clf-inf" and version="1.1"
INTO OUTFILE '/var/lib/mysql-files/logTime.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n';


SELECT 'epoch', 'image_num', 'runtime', 'pred_total', 'pred_transfer', 'pred_deploy', 'pred_proc', 'act_total', 'act_transfer', 'act_deploy', 'act_proc'
UNION
SELECT UNIX_TIMESTAMP(time_stamp) as epoch, image_num, runtime, pred_total, pred_transfer, pred_deploy, pred_proc,
act_total, act_transfer, act_deploy, act_proc from LogTime where task_id >= 512 and task_id < 1343
INTO OUTFILE '/var/lib/mysql-files/logTime_4runtimes.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n';


SELECT 'epoch', 'image_num', 'runtime', 'pred_total', 'pred_transfer', 'pred_deploy', 'pred_proc', 'act_total', 'act_transfer', 'act_deploy', 'act_proc'
UNION
SELECT UNIX_TIMESTAMP(time_stamp) as epoch, image_num, runtime, pred_total, pred_transfer, pred_deploy, pred_proc,
act_total, act_transfer, act_deploy, act_proc from LogTime where task_id >= 6457
INTO OUTFILE '/var/lib/mysql-files/logTime_4runtimes.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n';

select * from LogTime where task_id >= 512 and task_id < 1343;
select * from LogTime where task_id > 1350;

--step-up experiment
select * from LogTime where task_id >= 6161;



SELECT 'epoch', 'cpu', 'gpu1', 'gpu2'
UNION
SELECT UNIX_TIMESTAMP(time_stamp) as epoch, cpu, gpu1, gpu2 from DeploymentTime where deployment_id > 2429
INTO OUTFILE '/var/lib/mysql-files/deployTime.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n';
