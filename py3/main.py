import os

import cv2
from flask import Response as Res

import sift
import pickle
import threading
import time
import uuid


# 服务器对http请求的响应
class Response:
    def __init__(self, code, msg, data):
        self.code = code
        self.msg = msg
        self.data = data

    @staticmethod
    def of_ok(data):
        return Response('OK', None, data)

    @staticmethod
    def of_fail(msg):
        return Response('FAIL', msg, None)

    def to_string(self):
        return self.__dict__


# 图片模组
class ImageModel:
    def __init__(self, image):
        self.image = image
        self.des = None  # 特征点
        self.word_summary = None  # 关键词直方图
        self.idf_word_summary = None  # idf倒排后关键词直方图


class Cache:
    def __init__(self):
        self.data: [ImageModel] = []
        self.idf_count = 0  # idf的数量
        self.word_list = None  # 关键词列表
        self.idf = None

    def add(self, img: ImageModel):
        self.data.append(img)


cache = Cache()


# 获取已上传的图片数量
def get_tot_cnt():
    return Response.of_ok(len(cache.data)).to_string()


# 获取cache中索引图片的数量
def get_cnt():
    return Response.of_ok(cache.idf_count).to_string()


def get_file(filename):
    with open('train/' + filename, 'rb') as f:
        result = Res(f.read(), mimetype="image/jpeg")
        print(result)
        return result


# 保存图片到cache中
def upload_file(filename):
    # fileArrays = files.split('\n')
    cache.add(ImageModel(filename))


# 检索图片
def reduce():
    start_time = time.time()
    global cache
    word_count = 0
    path = "./train"
    for lists in os.listdir(path):
        sub_path = os.path.join(path, lists)
        if os.path.isfile(sub_path):
            word_count = word_count + 1

    # 创建线程
    class ThreadRunner(threading.Thread):
        def __init__(self, function, index):
            super().__init__()
            self.function = function
            self.index = index

        def run(self) -> None:
            self.function(self.index)

        @staticmethod
        def add(function, count):
            thread_list = []
            for i in range(count):
                cur_thread = ThreadRunner(function, i)
                cur_thread.start()
                thread_list.append(cur_thread)
            for thread in thread_list:
                thread.join()

    def get_des(index):
        if cache.data[index].des is None:
            cache.data[index].des = sift.get_des('./train/' + cache.data[index].image)

    ThreadRunner.add(get_des, len(cache.data))
    cache.word_list = sift.encode([data.des for data in cache.data], word_count)
    for img in cache.data:
        img.word_summary = sift.get_word_summary(img.des, cache.word_list)
    cache.idf = sift.tf_idf([data.word_summary for data in cache.data])
    for img in cache.data:
        img.idf_word_summary = sift.idf_render(img.word_summary, cache.idf)

    cache.idf_count = len(cache.data)
    end_time = time.time()
    print("retrieve {} image cost: {}s".format(cache.idf_count, end_time - start_time))
    with open(str(uuid.uuid4()), 'wb') as f:
        pickle.dump(cache, f)  # 结果存储到文件
    return Response.of_ok('OK').to_string()


# 图像查找
def find_image(filepath):
    image = cv2.imread(filepath)
    img_count = 10
    # image.save('./test/' + image.filename)
    des = sift.get_des('./test/' + image.filename)
    word_summary = sift.get_word_summary(des, cache.word_list)
    idf_word_summary = sift.idf_render(word_summary, cache.idf)
    res = []
    image_label = image.filename.split('_')[0]
    total_cnt = 0
    top5_cnt = 0
    top_10_cnt = 0
    p50_cnt = 0
    p50_ac_cnt = 0
    for i in range(cache.idf_count):
        match_value = sift.summary_match(idf_word_summary, cache.data[i].idf_word_summary)
        if match_value > 0.5:
            p50_cnt += 1
            if cache.data[i].image.split('_')[0] == image_label:
                total_cnt += 1
                p50_ac_cnt += 1
        elif cache.data[i].image.split('_')[0] == image_label:
            total_cnt += 1
        res.append((match_value, cache.data[i].image))
    res.sort(key=lambda x: x[0], reverse=True)
    for i in range(5):
        if res[i][1].split('_')[0] == image_label:
            top5_cnt += 1
            top_10_cnt += 1
    for i in range(5, 10):
        if res[i][1].split('_')[0] == image_label:
            top_10_cnt += 1
    return Response.of_ok({'info': {'tp5_ac': "{:.2f}%".format(top5_cnt / 5 * 100),
                                    'tp5_rc': "{:.2f}%".format(top5_cnt / total_cnt * 100),
                                    'tp10_ac': "{:.2f}%".format(top_10_cnt / 10 * 100),
                                    'tp10_rc': "{:.2f}%".format(top_10_cnt / total_cnt * 100),
                                    'p50_ac': "{:.2f}%".format(p50_ac_cnt / p50_cnt * 100),
                                    'p50_rc': "{:.2f}%".format(p50_ac_cnt / total_cnt * 100)
                                    },
                           'data': [{'value': img[0], 'name': img[1]} for img in res[:img_count]]}).to_string()
