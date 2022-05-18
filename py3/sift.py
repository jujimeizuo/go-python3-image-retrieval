import cv2
import math
import numpy as np
from scipy.cluster.vq import kmeans, vq

# 实例化sift函数
sift = cv2.SIFT_create()


# 获取特征点和特征向量
def get_des(image_path):
    image = cv2.imread(image_path)
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)  # 灰度直方图
    kp, des = sift.detectAndCompute(gray, None)  # kp is 关键点， des is sift特征向量

    ''' 画出关键点
    print(des)
    ret = cv2.drawKeypoint's(gray, kp, image)
    cv2.imshow('ret', ret)
    cv2.waitKey(0)
    cv2.destroyAllWindows()
    '''

    return des


# kmeans聚类，获取关键词
def encode(image_des_list, word_num):
    image_des_stack = None
    for des in image_des_list:
        if image_des_stack is None:
            image_des_stack = des
        else:
            image_des_stack = np.vstack((image_des_stack, des))
    word_list, distortion = kmeans(image_des_stack, word_num)  # 欧式距离聚类
    return word_list


# 将特征向量改为关键词直方图(TF)
def get_word_summary(image_des, word_list):
    word_summary = np.zeros(len(word_list))
    own_words, distance = vq(image_des, word_list)
    for own_word in own_words:
        word_summary[own_word] += 1
    return word_summary


# 通过IDF进行关键词倒排
def tf_idf(image_word_summary_list):
    idf = np.zeros(len(image_word_summary_list[0]))
    for i in range(len(idf)):
        for image_word_summary in image_word_summary_list:
            idf[i] += (1 if image_word_summary[i] > 0 else 0)
        idf[i] = math.log(len(image_word_summary_list) / (idf[i] + 1))
    return idf


# 获得加权关键词直方图
def idf_render(image_word_summary, idf):
    image_word_summary *= idf
    image_word_summary = image_word_summary / np.linalg.norm(image_word_summary)
    return image_word_summary


# 通过余弦相似度，获得距离最近的图片
def summary_match(summary1, summary2):
    value = float(np.dot(summary1, summary2))  # 计算特征之间的余弦距离（相似度）
    return value
