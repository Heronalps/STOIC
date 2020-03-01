import re
import pandas as pd

runtime_re = re.compile("The task is scheduled at (.*)")
imagenum_re = re.compile("Start running task image-clf-inf version 2.2 on (\d*) images")
designated_runtimes = []
image_nums = []
with open("./client_log.txt", "r") as file:
    for line in file:
        runtime_match = runtime_re.search(line)
        imagenum_match = imagenum_re.search(line)
        if runtime_match != None:
            runtime = runtime_match.group(1)
            designated_runtimes.append(runtime.strip())
        if imagenum_match != None:
            imagenum = imagenum_match.group(1)
            image_nums.append(imagenum)

print(designated_runtimes)
print(image_nums)
print (len(designated_runtimes))
print ( len(image_nums))