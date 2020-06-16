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
    for t in range(1, 9):
        with Pool(t) as p:
            print ("num of threads : " + str(t))
            arr = []
            for i in range(t):
                arr.append(1024)
            print(p.map(f, arr))