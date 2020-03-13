import numpy as np
import pandas as pd

df_logtime4 = pd.read_csv("./logTime_4runtimes.csv")

counter = 0
for index in range(2, len(df_logtime4), 4):
    if df_logtime4[index:index+1]["act_total"].item() < df_logtime4[index+1:index+2]["act_total"].item():
            counter += 1
print ("Speed up run: ", counter)
print ("Speed up rate : ", counter / (len(df_logtime4)/4))

df_gpu = pd.read_csv("./logTime_gpu_2.csv")

counter = 0
for index in range(0, len(df_gpu), 2):
    if df_gpu[index:index+1]["act_proc"].item() > df_gpu[index+1:index+2]["act_proc"].item():
        # print (df_gpu[index:+index+1])
        # print (df_gpu[index+1:+index+2])
        counter += 1
print ("Speed up count : ", counter)
print ("Speed up rate : ", round(counter/(len(df_gpu)/2), 2))
