#!/usr/bin/python
# -*- coding: UTF-8 -*-
import copy
import io

from fontTools import subset
from fontTools.ttLib import TTFont


def process():
    font = TTFont('./SourceHanSansCN-Regular.otf')

    file = io.open("./high_usage_chinese", mode='r', encoding='UTF-8')
    highUsageChinese = file.read()

    subSetter = subset.Subsetter()
    subSetter.populate(text=highUsageChinese)
    subSetter.subset(font)

    cmap = font.getBestCmap()

    tempMap = copy.deepcopy(cmap)
    for k, v in tempMap.items():
        cmap[k + 10] = cmap.pop(k)

    font.save('./font.woff2')


if __name__ == '__main__':
    process()
