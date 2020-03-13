SELECT 'epoch', 'ts', 'cpu', 'gpu1', 'gpu2'
UNION
SELECT UNIX_TIMESTAMP(time_stamp) as epoch, time_stamp, cpu, gpu1, gpu2 from DeploymentTime INTO OUTFILE '/var/lib/mysql-files/deploymentTime.csv' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n';