import re
import pandas as pd

regex = re.compile("The task is scheduled at (.*)")
designated_runtimes = []
with open("./client_log.txt", "r") as file:
    for line in file:
        result = regex.search(line)
        if result != None:
            runtime = result.group(1)
            designated_runtimes.append(runtime)

df = pd.Series(data=designated_runtimes)
print(df)