import sift
import pickle
import threading
import time
import uuid
import numpy as np


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

        self.class_word = None
        self.class_idf = None
        self.class_idf_word_summary = None

    def add(self, img: ImageModel):
        self.data.append(img)


cache = Cache()


# 获取已上传的图片数量
def get_tot_cnt():
    return Response.of_ok(len(cache.data)).to_string()


# 获取cache中索引图片的数量
def get_cnt():
    return Response.of_ok(cache.idf_count).to_string()


# 保存图片到cache中
def upload_file(filename):
    cache.add(ImageModel(filename))


# 检索图片
def reduce_tmp():
    start_time = time.time()
    global cache
    word_count = 1024

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
            cache.data[index].des = sift.get_des('./img/' + cache.data[index].image)

    ThreadRunner.add(get_des, len(cache.data))
    cache.word_list = sift.encode_plus([data.des for data in cache.data], word_count)
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


def reduce():
    global cache
    class_des = dict()
    class_word = 1024
    for image in cache.data:
        class_name = image.image.split('_')[0]
        if class_name not in class_des.keys():
            class_des[class_name] = image.idf_word_summary
        else:
            class_des[class_name] = np.vstack((class_des[class_name], image.idf_word_summary))
    cache.class_word = sift.encode(class_des.values(), class_word)
    class_word_summary = dict()
    for key, value in class_des.items():
        class_word_summary[key] = sift.get_word_summary(value, cache.class_word)
    cache.class_idf = sift.tf_idf(list(class_word_summary.values()))
    cache.class_idf_word_summary = dict()
    for key, value in class_word_summary.items():
        cache.class_idf_word_summary[key] = sift.idf_render(value, cache.class_idf)
    return Response.of_ok('OK').to_string()


def find_image(filename):
    img_count = 10
    des = sift.get_des('./img/' + filename)
    word_summary = sift.get_word_summary(des, cache.word_list)
    idf_word_summary = sift.idf_render(word_summary, cache.idf)

    res = []

    # region 重排序
    best_match = dict()
    # endregion

    image_label = filename.split('_')[0]
    total_cnt = 0
    top5_cnt = 0
    top10_cnt = 0
    p50_cnt = 0
    p50_ac_cnt = 0
    for i in range(cache.idf_count):
        match_value = sift.summary_match(idf_word_summary, cache.data[i].idf_word_summary)
        cur_name = cache.data[i].image.split('_')[0]

        # region 重排序
        if i != 0:
            if cur_name not in best_match.keys():
                best_match[cur_name] = {
                    'cnt': 1,
                    'value': match_value
                }
            elif best_match[cur_name]['cnt'] < 3:
                best_match[cur_name]['cnt'] += 1
                best_match[cur_name]['value'] += match_value
        # endregion

        if match_value > 0.5:
            p50_cnt += 1
            if cur_name == image_label:
                total_cnt += 1
                p50_ac_cnt += 1
        elif cur_name == image_label:
            total_cnt += 1
        res.append((match_value, cache.data[i].image))
    for key, value in best_match.items():
        best_match[key] = value['value'] / value['cnt']
    res = [(value * best_match[name.split('_')[0]], name) for value, name in res]
    res.sort(key=lambda x: x[0], reverse=True)

    ap_value = 0
    ap_cnt = 0
    for i in range(5):
        if res[i][1].split('_')[0] == image_label:
            top5_cnt += 1

    for i in range(10):
        if res[i][1].split('_')[0] == image_label:
            top10_cnt += 1

    for i in range(len(res)):
        if res[i][1].split('_')[0] == image_label:
            ap_cnt += 1
            ap_value += ap_cnt / (i + 1)

    return Response.of_ok({'info': {'tp5_ac': "{:.2f}%".format(top5_cnt / 5 * 100),
                                    'tp5_rc': "{:.2f}%".format(top5_cnt / total_cnt * 100),
                                    'tp10_ac': "{:.2f}%".format(top10_cnt / 10 * 100),
                                    'tp10_rc': "{:.2f}%".format(top10_cnt / total_cnt * 100),
                                    'p50_ac': "{:.2f}%".format(p50_ac_cnt / p50_cnt * 100),
                                    'p50_rc': "{:.2f}%".format(p50_ac_cnt / total_cnt * 100),
                                    'ap': ap_value / ap_cnt
                                    },
                           'data': [{'value': img[0], 'name': img[1]} for img in res[:img_count]]}).to_string()


# 图像查找
def find_image_tmp(filename):
    img_count = 10
    des = sift.get_des('./img/' + filename)
    word_summary = sift.get_word_summary(des, cache.word_list)
    idf_word_summary = sift.idf_render(word_summary, cache.idf)
    res = []
    image_label = filename.split('_')[0]
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
