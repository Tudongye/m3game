# -*- coding: utf-8 -*-
#!/usr/bin/python3
# -*- coding: utf-8 -*-
# 用于分析目录内go package的引用关系
# python analysis_golang_package.py <目录名>
import os,re,sys

go_path = sys.argv[1]

prefix = '''digraph go_pkg_relation {
    graph [
        rankdir = "LR"
        //splines=polyline
        overlap=false
    ];

    node [
        fontsize = "16"
        shape = "ellipse"\r
    ];

    edge [
    ];
'''

suffix = '}'

#从当前文件中获取包名
def get_package_of_file(path_file):
    with open(path_file, 'r', encoding="utf-8") as file_input:
        tmpline = file_input.readline()
        while (tmpline):
            #匹配 package xxxxxx
            m = re.search(r'^package\s*([0-9a-zA-Z_\-]*)', tmpline)
            if m:
                return m.group(1)
            tmpline = file_input.readline()

#从当前文件中获取导入的包
def get_imported_package_from_file(path_file):
    imported_package = []
    with open(path_file, 'r', encoding="utf-8") as file_input:
        tmpline = file_input.readline()
        in_import = False
        while (tmpline):
            #单行import的处理
            m = re.search(r'^import\s*\"([0-9a-zA-Z/\.])\"', tmpline)
            if m:
                imported_package.append(m.group(1))

            #多行import的处理
            m = re.search(r'^import\s*\(', tmpline)
            if m:
                in_import = True
                tmpline = file_input.readline()
                continue
            if in_import:
                m = re.search(r'^\s*\)', tmpline)
                if m:
                    in_import = False
                    tmpline = file_input.readline()
                    continue
                #多行import，最后的)不换行
                m = re.search(r'\"(.*/)([0-9a-zA-Z]*)\"\s*\)', tmpline)
                if m:
                    imported_package.append(m.group(2))
                    in_import = False
                    tmpline = file_input.readline()
                    continue
                m = re.search(r'\"(.*/)([0-9a-zA-Z]*)\"', tmpline)
                if m:
                    imported_package.append(m.group(2))
            tmpline = file_input.readline()
    return imported_package

def build_node_from_file(file_path, file_name):
    i = 0
    space4 = '    '
    space8 = space4 + space4
    space12 = space4 + space8
    space16 = space4 + space12
    
    #获取文件的package和导入的package
    file_package = get_package_of_file(os.path.join(file_path, file_name))
    imported_packages = get_imported_package_from_file(os.path.join(file_path, file_name))
    #node_str=
    #    "fweh" [
    #        label = "<head> fweh.h\l|
    #            {|{
    #                <unaligned> unaligned.h\l|
    #                <skbuff> skbuff.h\l|
    #                <if_ether> if_ether.h\l|
    #                <if> if.h\l|
    #            }}"
    #        shape = "record"
    #    ];
    if not file_package:
        return
    if file_package not in node_created_db:
        node_created_db[file_package] = {}
        node_created_db[file_package]["imported_package"] = {}
        node_created_db[file_package]["node_str_mid"] = ""
        node_created_db[file_package]["node_str_pre"] = space4 + '\"' + file_package + '\" [\n' + space8 + 'label = \"<head> '+ file_package +'\l|\n' + space12 + '{|{\n'
        node_created_db[file_package]["node_str_suffix"] = space12 + '}}\"\n' + space8 + 'shape = \"record\"\n' + space4 + '];\n'



    for imported_package in imported_packages:
        #导入的包尚未被记录过，记录到db中，
        if imported_package not in node_created_db[file_package]["imported_package"]:
            node_created_db[file_package]["node_str_mid"] = node_created_db[file_package]["node_str_mid"] + space16 + '<'+ imported_package + '> ' + imported_package +  '\l|\n'
            #记录当前包指向被包含包的边
            tmp_edge_str = space4 + file_package + ':' + imported_package + ' -> ' + imported_package + ':' + 'head' #+ '[color="' + color_arrary[i%len(color_arrary)] + '"]\n'
            node_created_db[file_package]["imported_package"][imported_package] = tmp_edge_str

#字典类型，记录已经创建的node信息，
#{
#    "package": #包名
#    {
#        "imported_package":{"包名":"边"}
#        "node_str_pre":"xxxxx", #描述node的字符串前缀
#        "node_str_mid":"xxxxx", #描述node的字符串中间部分
#        "node_str_suffix":"xxxxx", #描述node的字符串后缀
#        "edge":"xxxxx", #描述node的字符串后缀
#    }
# }
node_created_db = {}
node_str = ''
edge_string = ''
color_arrary = ['red', 'green', 'blue', 'black','blueviolet','brown', 'cadetblue','chocolate','crimson','cyan','darkgrey','deeppink','darkred']
i = 0

for maindir, subdir, file_name_list in os.walk(go_path):#('G:\git_repository\linux-stable\linux-4.18\drivers\usb'):
    for tmpfile in file_name_list:
        if re.match(r'.*\.go$', tmpfile):
            result = build_node_from_file(maindir, tmpfile)

for package in node_created_db:
    node_str = node_str + node_created_db[package]["node_str_pre"] + node_created_db[package]["node_str_mid"] + node_created_db[package]["node_str_suffix"]
    for imported_package in node_created_db[package]["imported_package"]:
        if imported_package in node_created_db:
            edge_string = edge_string + node_created_db[package]["imported_package"][imported_package] + '[color="' + color_arrary[i%len(color_arrary)] + '"]\n'
            i = i + 1

print(prefix + node_str + '\n' + edge_string + suffix)
