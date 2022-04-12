#!/usr/bin/python3
# -*- coding: UTF-8 -*-
import io

from fontTools import subset
from fontTools.ttLib import TTFont
import xml.etree.ElementTree as ET


class ObfsFont:
    def __init__(self):
        self.origin_font_file = './SourceHanSansCN-Regular.otf'
        self.font_file = 'font.woff2'

    def process(self):
        temp_xml_path = './temp.xml'

        file = io.open("./high_usage_chinese", mode='r', encoding='UTF-8')
        high_usage_chinese = file.read()

        font = TTFont(self.origin_font_file)

        sub_setter = subset.Subsetter()
        sub_setter.populate(text=high_usage_chinese)
        sub_setter.subset(font)

        font.saveXML(temp_xml_path)
        tree = ET.parse(temp_xml_path)
        for pcamp in tree.getroot().iter('cmap'):
            for scmap in pcamp.iter('cmap_format_4'):
                for item in scmap.iter('map'):
                    code = item.get('code')
                    item.set('code', hex(int(code, 16) + 10))

        tree.write(temp_xml_path)

        font = TTFont()
        font.importXML(fileOrPath='./temp2.xml')
        font.save(self.font_file)

        os.remove(temp_xml_path)


if __name__ == '__main__':
    font_processor = ObfsFont()
    font_processor.process()
