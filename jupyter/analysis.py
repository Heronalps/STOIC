import numpy as np
import pandas as pd
import sys, os

df_logtime4 = pd.read_csv("./logTime_4runtimes.csv")
df_stoic = pd.read_csv("./logTime_stoic.csv")

df_edge = df_logtime4[df_logtime4["runtime"]=="edge"]
for index, row in df_edge.iterrows():
    if index % 4 != 0:
        print (index)
        sys.exit(0)

df_edge.reset_index(inplace=True)

df_stoic_front = df_stoic[0:len(df_edge)]
df_edge[df_stoic_front["image_num"] != df_edge["image_num"]]

runtimes = []
for index in range(0, len(df_logtime4), 4):
    minTime = float("inf")
    runtime = ""
    for i in range(index, index + 4):
        if (df_logtime4[i: i+1]["act_total"].item() < minTime):
            minTime = df_logtime4[i: i + 1]["act_total"].item()
            runtime = df_logtime4[i: i + 1]["runtime"].item()
    runtimes.append(runtime)

groundtruth_runtimes = pd.Series(runtimes)
stoic_runtimes = df_stoic_front["runtime"]
positive = stoic_runtimes[stoic_runtimes == groundtruth_runtimes]

print("STOIC True Positive Rate: {0}".format(len(positive) / len(stoic_runtimes)))

mae = sum(abs(df_stoic["pred_total"] - df_stoic["act_total"])) / len(df_stoic)
print ("STOIC Mean Absolute Error: {0}".format(mae))