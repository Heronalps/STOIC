from multiprocessing import Process, Queue
import numpy as np
import time, math, os, argparse, shutil
from gpuinfo import GPUInfo
from zipfile import ZipFile


class_list = ["Birds", "Empty", "Fox", "Humans", "Rodents"]
MODEL_DIR = '/racelab/checkpoints/resnet50_model.h5'
IMG_DIR = '/racelab/SantaCruzIsland_Labeled_5Class/Birds'
PROBE_IMG = '/racelab/SantaCruzIsland_Labeled_5Class/Birds/IMG_0198.JPG'
TARGET_IMG = '/IMG_0198.JPG'
WIDTH = 1920
HEIGHT = 1080

class Scheduler:
    def __init__(self, gpu_num):
        self._queue = Queue()
        self._gpu_num = gpu_num
        self.__init_workers()

    def __init_workers(self):
        self._workers = list()
        for gpuid in range (self._gpu_num):
            self._workers.append(Worker(gpuid, self._queue))

    def start(self, image_list):
        for img in image_list:
            self._queue.put(img)

        # Add a None to indicate the end of queue
        self._queue.put(None)

        for worker in self._workers:
            worker.start()

        for worker in self._workers:
            worker.join()
        print ("All image are done inferencing...")

class Worker(Process):
    def __init__(self, gpuid, queue):
        Process.__init__(self, name="ModelProcessor")
        self._gpuid = gpuid
        self._queue = queue
    
    def run(self):
        #set enviornment
        os.environ["CUDA_DEVICE_ORDER"] = "PCI_BUS_ID"
        os.environ["CUDA_VISIBLE_DEVICES"] = str(self._gpuid)

        from tensorflow.keras.applications.resnet50 import preprocess_input
        from tensorflow.keras.preprocessing import image
        from tensorflow.keras.models import load_model
        trained_model = load_model(MODEL_DIR)
        
        while True:
            img_path = self._queue.get()
            if img_path == None:
                self._queue.put(None)
                break
            img = image.load_img(path=img_path, target_size=(1920, 1080))
            x = image.img_to_array(img)
            x = np.expand_dims(x, axis=0)
            x = preprocess_input(x)
            y_prob = trained_model.predict(x)
            index = y_prob.argmax()
            print ("image : {0}, index : {1}".format(img_path, index))
        
        print("GPU {} has done inferencing...".format(self._gpuid))

# For the runtime with 0 GPU
def run_sequential(image_list):
    from tensorflow.keras.applications.resnet50 import preprocess_input
    from tensorflow.keras.preprocessing import image
    from tensorflow.keras.models import load_model

    trained_model = load_model(MODEL_DIR)
    for img_path in image_list:
        img = image.load_img(path=img_path, target_size=(1920, 1080))
        x = image.img_to_array(img)
        x = np.expand_dims(x, axis=0)
        x = preprocess_input(x)
        y_prob = trained_model.predict(x)
        index = y_prob.argmax()
        print ("image : {0}, index : {1}".format(img_path, index))
    

def handler(event, context): 
    TEMP_DIR = os.getcwd() + "/image_buffer"
    os.mkdir(TEMP_DIR)
    zip_flag = False
    
    if isinstance(event['data'], dict) and "zip_path" in event['data']:
        global ZIP_PATH
        ZIP_PATH = event['data']['zip_path']
        zip_flag = True
        with ZipFile(ZIP_PATH, 'r') as zipFile:
            zipFile.extractall(TEMP_DIR)
    else:
        shutil.copyfile(PROBE_IMG, TEMP_DIR + TARGET_IMG)

    # Get GPU counts
    NUM_GPU = 0
    available_devices = GPUInfo.check_empty()
    if available_devices != None:
        NUM_GPU = len(available_devices)
    print ("Current GPU num is {0}".format(NUM_GPU))
    
    # Get image number
    num_image = 0
    image_list = list()
    for img in os.listdir(TEMP_DIR):
        image_list.append(os.path.join(TEMP_DIR, img))
        num_image += 1
        
    start = time.time()

    if NUM_GPU == 0:
        run_sequential(image_list)
    else:
        # initialize Scheduler
        scheduler = Scheduler(NUM_GPU)
        # start multiprocessing
        scheduler.start(image_list)
        
    end = time.time()

    # Clean up temp image folder
    shutil.rmtree(TEMP_DIR)
    if (zip_flag):
        os.remove(ZIP_PATH)

    print ("Time with model loading {0} for {1} images.".format(end - start, num_image))
    return ("Time with model loading {0} for {1} images.".format(end - start, num_image))

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('-p', '--zip_path')
    args = parser.parse_args()
    if args.zip_path:
        handler({"data" : {"zip_path" : args.zip_path}}, {})
    else:
        handler({"data" : {}}, {})