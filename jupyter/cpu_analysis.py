import re
import pandas as pd

regex = re.compile("nautilus\.io\/processor: (.*)")
cpu_dict = dict()
counter = 0
with open("./nodes.txt", "r") as file:
    for line in file:
        result = regex.search(line)
        if result != None:
            counter += 1
            cpu_type = result.group(1)
            if cpu_type in cpu_dict != None:
                cpu_dict[cpu_type] += 1
            else:
                cpu_dict[cpu_type] = 1
print ("cpu types: ", len(cpu_dict))
print ("nodes : ", counter)
cpu_dict_sorted = {k:v for k, v in sorted(cpu_dict.items(), key = lambda item: item[1], reverse=True)}

df = pd.DataFrame.from_dict(data=cpu_dict_sorted, columns=['num'], orient="index")
print(df)