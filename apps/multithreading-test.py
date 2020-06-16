from multiprocessing import Pool
import time

def f(x):
    start = time.time()
    while (True):
        curr = time.time()
        if (curr - start >= 300):
            print ("Epoch : " + str(time.ctime(curr)))
            break
        r = pow(x, x)

if __name__ == '__main__':
    for i in range(1, 9):
        with Pool(i) as p:
            print(p.map(f, [1024]))