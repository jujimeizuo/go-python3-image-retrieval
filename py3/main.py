import threading
import time

import sift


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


# 保存图片到cache中
def upload_file(filename):
    cache.add(ImageModel(filename))


def reduce():
    start_time = time.time()
    global cache
    word_count = 1024

    class TreadRunner(threading.Thread):
        def __init__(self, function, index):
            super().__init__()
            self.function = function
            self.index = index

        def run(self) -> None:
            self.function(self.index)

        @staticmethod
        def aha(function, count):
            thread_list = []
            for i in range(count):
                cur_thread = TreadRunner(function, i)
                cur_thread.start()
                thread_list.append(cur_thread)
            for thread in thread_list:
                thread.join()

    def get_des(index):
        if cache.data[index].des is None:
            cache.data[index].des = sift.get_des('image/' + cache.data[index].image)

    TreadRunner.aha(get_des, len(cache.data))
    print("get_des OK")
    cache.word_list = sift.encode([mod.des for mod in cache.data], word_count)
    print("encode OK")
    for img in cache.data:
        img.word_summary = sift.get_word_summary(img.des, cache.word_list)
    print("get_word_summary OK")
    cache.idf = sift.tf_idf([mod.word_summary for mod in cache.data])
    print("tf_idf OK")
    for img in cache.data:
        img.idf_word_summary = sift.idf_render(img.word_summary, cache.idf)
    print("idf_render OK")
    cache.idf_count = len(cache.data)
    end_time = time.time()
    print("reduce OK")
    print("reduce {} image cost: {}s".format(cache.idf_count, end_time - start_time))

    return Response.of_ok('Ok').to_string()


def find_image(filename):
    img_count = 10
    des = sift.get_des('./img/' + filename)
    word_summary = sift.get_word_summary(des, cache.word_list)
    idf_word_summary_first = sift.idf_render(word_summary, cache.idf)

    res = []

    for i in range(cache.idf_count):
        match_value = sift.summary_match(idf_word_summary_first, cache.data[i].idf_word_summary)
        res.append((match_value, cache.data[i].image))
    res.sort(key=lambda x: x[0], reverse=True)

    print("resort!!!")
    cache_temp = Cache()
    for i in range(30):
        cache_temp.add(ImageModel(res[i][1]))

    start_time = time.time()
    word_count = 1024

    class TreadRunner(threading.Thread):
        def __init__(self, function, index):
            super().__init__()
            self.function = function
            self.index = index

        def run(self) -> None:
            self.function(self.index)

        @staticmethod
        def aha(function, count):
            thread_list = []
            for index in range(count):
                cur_thread = TreadRunner(function, index)
                cur_thread.start()
                thread_list.append(cur_thread)
            for thread in thread_list:
                thread.join()

    def get_des(index):
        if cache_temp.data[index].des is None:
            cache_temp.data[index].des = sift.get_des('image/' + cache_temp.data[index].image)

    TreadRunner.aha(get_des, len(cache_temp.data))
    print("get_des OK")
    cache_temp.word_list = sift.encode([mod.des for mod in cache_temp.data], word_count)
    print("encode OK")
    for img in cache_temp.data:
        img.word_summary = sift.get_word_summary(img.des, cache_temp.word_list)
    print("get_word_summary OK")
    cache_temp.idf = sift.tf_idf([mod.word_summary for mod in cache_temp.data])
    print("tf_idf OK")
    for img in cache_temp.data:
        img.idf_word_summary = sift.idf_render(img.word_summary, cache_temp.idf)
    print("idf_render OK")
    cache_temp.idf_count = len(cache_temp.data)
    end_time = time.time()
    print("reduce OK")
    print("reduce {} image cost: {}s".format(cache_temp.idf_count, end_time - start_time))

    word_summary = sift.get_word_summary(des, cache_temp.word_list)
    idf_word_summary_second = sift.idf_render(word_summary, cache_temp.idf)

    res = []
    best_match = dict()

    image_label = filename.split('_')[0]
    total_cnt = len(cache_temp.data)
    top5_cnt = 0
    top10_cnt = 0
    p50_cnt = 0
    p50_ac_cnt = 0
    for i in range(cache_temp.idf_count):
        match_value = sift.summary_match(idf_word_summary_second, cache_temp.data[i].idf_word_summary)
        cur_name = cache_temp.data[i].image.split('_')[0]

        # region 扩展查询
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
                p50_ac_cnt += 1
        res.append((match_value, cache_temp.data[i].image))
    for key, value in best_match.items():
        best_match[key] = value['value'] / value['cnt']
    res = [(value * best_match[name.split('_')[0]], name) for value, name in res]
    res.sort(key=lambda x: x[0], reverse=True)

    for i in range(5):
        if res[i][1].split('_')[0] == image_label:
            top5_cnt += 1

    for i in range(10):
        if res[i][1].split('_')[0] == image_label:
            top10_cnt += 1

    return Response.of_ok({'info': {'tp5_ac': "{:.2f}%".format(top5_cnt / 5 * 100),
                                    'tp5_rc': "{:.2f}%".format(top5_cnt / total_cnt * 100),
                                    'tp10_ac': "{:.2f}%".format(top10_cnt / 10 * 100),
                                    'tp10_rc': "{:.2f}%".format(top10_cnt / total_cnt * 100),
                                    'p50_ac': "{:.2f}%".format(p50_ac_cnt / p50_cnt * 100),
                                    'p50_rc': "{:.2f}%".format(p50_ac_cnt / total_cnt * 100),
                                    },
                           'data': [{'value': img[0], 'name': img[1]} for img in res[:img_count]]}).to_string()
