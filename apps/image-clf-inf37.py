from tensorflow.keras.applications.resnet50 import preprocess_input
from tensorflow.keras.preprocessing import image
from tensorflow.keras.models import load_model
from tensorflow.python.client import device_lib
import numpy as np
import time, math, os
import argparse

class_list = ["Birds", "Empty", "Fox", "Humans", "Rodents"]
NUM_IMAGE = 1
INF_DIR = "/racelab/SantaCruzIsland_Labeled_5Class/Birds"
WIDTH = 1920
HEIGHT = 1080

# PATH = "/racelab/data/SantaCruzIsland_Validation_5Class/Birds"


def handler(event, context): 
    start1 = time.time()
    # if isinstance(event['data'], dict) and "path" in event['data']:
    #     global PATH
    #     PATH = event['data']['path']

    if isinstance(event['data'], dict) and "num_image" in event['data']:
        global NUM_IMAGE
        NUM_IMAGE = int(event['data']['num_image'])
    
    # Parallel with multiple GPUs
    available_devices = device_lib.list_local_devices()
    NUM_GPU = len([x for x in available_devices if x.device_type == 'GPU'])
    print ("Current GPU num is {0}".format(NUM_GPU))

    # Increase BATCH_SIZE based on number of GPUs to harness the quasi-linear speedup of multiple GPUS
    # Each GPU takes 2 augmented images for training at one epoch
    
    # BATCH_SIZE = 1
    # inf_datagen = image.ImageDataGenerator(preprocessing_function=preprocess_input, rotation_range=90, \
    #                                        horizontal_flip=True, vertical_flip=True)
    # inf_generator = inf_datagen.flow_from_directory(INF_DIR, target_size=(WIDTH, HEIGHT), \
    #                                                 batch_size = BATCH_SIZE)
    
    BATCH_SIZE = 2 * NUM_GPU if NUM_GPU > 0 else 2
    trained_model = load_model('/racelab/checkpoints/resnet50_model.h5')

    start2 = time.time()
    counter = 0
    for img in os.listdir(INF_DIR):
        PATH = INF_DIR + '/' + img
        img = image.load_img(path=PATH, target_size=(1920, 1080))
        x = image.img_to_array(img)
        x = np.expand_dims(x, axis=0)
        x = preprocess_input(x)
        y_prob = trained_model.predict(x)
        index = y_prob.argmax()
        counter += 1
        if counter == NUM_IMAGE:
            break
        print ("image : {0}, index : {1}".format(PATH, index))

    # y_pred = trained_model.predict_generator(inf_generator, steps = math.ceil(NUM_IMAGE / BATCH_SIZE), workers=8)
    # print ("shape : ", y_pred.shape)
    # index = y_prob.argmax()
    # print ("index : ", index)

    # return ("Predicted : {0}  Probability : {1} Time with model loading: {2} Time without model loading {3} for {4} images.".format(class_list[index], y_prob[0][index], (time.time() - start1), (time.time() - start2), NUM_IMAGE))
    print ("Time with model loading: {0} Time without model loading {1} for {2} images.".format(time.time() - start1, time.time() - start2, NUM_IMAGE))
    return ("Time with model loading: {0} Time without model loading {1} for {2} images.".format(time.time() - start1, time.time() - start2, NUM_IMAGE))

if __name__ == "__main__":
    handler({"data" : {"num_image" : 10}}, {})