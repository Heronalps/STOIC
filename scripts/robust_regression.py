import pymysql.cursors
import argparse
import numpy as np
from sklearn import linear_model, datasets


def handler(runtime, app, version, numDP):
    # Connect to the database
    connection = pymysql.connect(host='localhost',
                                user='root',
                                password='Stoic!@#$%^123456',
                                db='test',
                                charset='utf8mb4',
                                cursorclass=pymysql.cursors.DictCursor)
    X = list()
    y = list()
    try:
        with connection.cursor() as cursor:
            # Read a single record
            sql = "SELECT `image_num`, `{0}` from `ProcessingTime{1}` WHERE `application` = %s " \
            "and `version` = %s ORDER BY `task_id` DESC LIMIT %s;".format(runtime, runtime.capitalize())
            cursor.execute(sql, (app, version, numDP,))
            all_rows = cursor.fetchall()
            for row in all_rows:
                X.append([row['image_num']])
                y.append([row[runtime]])
    finally:
        connection.close()                            
    X = np.array(X)
    y = np.array(y)
    
    ransac = linear_model.RANSACRegressor()
    ransac.fit(X, y)
    coef = ransac.estimator_.coef_
    intercept = ransac.estimator_.intercept_
    print ("Coefficient: {0} Intercept: {1}".format(coef[0][0], intercept[0]))
    return "Coefficient: {0} Intercept: {1}".format(coef[0][0], intercept[0])

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("runtime")
    parser.add_argument("app")
    parser.add_argument("version")
    parser.add_argument("numDP")
    args = parser.parse_args()
    handler(runtime = args.runtime, app = args.app, version = args.version, numDP = int(args.numDP))

# python robust_regression.py cpu image-clf-inf 1.7 10