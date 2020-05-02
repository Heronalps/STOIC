#! /usr/bin/env python3

import os
import sys
import json
import boto3
from zipfile import ZipFile

TEMP_DIR = os.getcwd() + "/image_buffer"

def main():
    s3 = boto3.resource('s3')
    bucket = s3.Bucket('seneca-racelab')

    bucket.download_file(Filename="temp.zip", Key='image_batch.zip')
    file_dir = os.getcwd() + "/temp.zip"

    with ZipFile(file_dir, 'r') as zipFile:
        zipFile.extractall(TEMP_DIR)

if __name__ == '__main__':
    main()