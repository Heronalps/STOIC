from multiprocessing import Process, Queue
import numpy as np
import time, math, os, argparse, shutil, pathlib, boto3
from gpuinfo import GPUInfo
from zipfile import ZipFile


class_list = ["Birds", "Empty", "Fox", "Humans", "Rodents"]
MODEL_DIR = '/racelab/checkpoints/resnet50_model.h5'
IMG_DIR = '/racelab/SantaCruzIsland_Labeled_5Class/Birds'
PROBE_IMG = '/racelab/SantaCruzIsland_Labeled_5Class/Birds/IMG_0198.JPG'
TEMP_DIR = "/racelab/image_buffer"
TARGET_IMG = '/IMG_0198.JPG'
WIDTH = 1920
HEIGHT = 1080

BUCKET = "test-bucket"
FILE_PATH = "/racelab/temp.zip"

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
    

def download_and_decompress_zip(key):
    s3 = boto3.resource('s3',
                        endpoint_url="http://rook-ceph-rgw-nautiluss3.rook",
                        aws_access_key_id="Y2WY7EJAPSVBC0PBP5IZ",
                        aws_secret_access_key= "MkNuKhdyoXPMsF3dHljnnLHOPp6KgGGxlurVwiMO")
    bucket = s3.Bucket(BUCKET)

    bucket.download_file(Key=key, Filename=FILE_PATH)

    with ZipFile(FILE_PATH, 'r') as zipFile:
        zipFile.extractall(TEMP_DIR)

def handler(event, context): 
    
    if event['data']['key_name'].lower() != "pseudo_key":
        key_name = event['data']['key_name']
        download_and_decompress_zip(key_name)
    else:
        pathlib.Path(TEMP_DIR).mkdir(parents=True, exist_ok=True)
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
        if img.lower().endswith(".jpg") or img.lower().endswith(".png"):
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
    

    print ("Time with model loading {0} for {1} images.".format(end - start, num_image))
    return ("Time with model loading {0} for {1} images.".format(end - start, num_image))

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('key_name')
    args = parser.parse_args()
    handler({"data" : {"key_name" : args.key_name}}, {})