import json
import xmind
import sys
import os
import sqlite3

def set_topic(root, val_data):
    root.setTitle(val_data['topic'])
    for sub_val in val_data['sub_topic'].values():
        set_topic(root.addSubTopic(), sub_val)


def set_one_sheet(root, input_file):
    data_dict = {}  # 将节点保存到字典中
    with open(input_file, 'r', encoding='utf-8') as f:
        json_data = json.load(f)
        for node in json_data['nodes']:
            data_dict[node['id']] = {
                'id': node['id'],
                'parentid': node['parentid'],
                'topic': node['topic'],
                'sub_topic': {},
            }

    data_root = None
    for val in data_dict.values():
        if val['id'] == 'root':
            data_root = val  # 保存根节点
        elif n := data_dict.get(val['parentid']):
            # 当前节点挂到父节点下
            n['sub_topic'][val['id']] = val
    if data_root is None:
        print('数据有误,没有找到根节点')
        return
    # print(json.dumps(data_root))  # 可以打印看看结果
    set_topic(root, data_root)  # 递归设置子节点


def convert_mindmap_xmind(input_file, save_file):
    if os.path.exists(save_file):
        os.remove(save_file)  # 目标文件存在则删除
    mind = xmind.load(save_file)

    sheet = None
    for title, path in input_file:
        if sheet:
            # 创建新的sheet
            sheet = mind.createSheet()
        else:
            # 创建默认sheet,打开后默认显示这个
            sheet = mind.getPrimarySheet()
        # 设置sheet标题
        sheet.setTitle(title)
        root = sheet.getRootTopic()
        # 设置属性为: 逻辑图(向右), 和有道云笔记保持一致
        root.setAttribute('structure-class', 'org.xmind.ui.logic.right')
        # 保存其中一个sheet
        set_one_sheet(root, path)
        print(title, '=>', path)
    # 将对象保存到文件中
    xmind.save(mind, save_file)

def get_mindmap_path(db, title):
    conn = sqlite3.connect(db)
    c = conn.cursor()
    # 传入title支持sql通配符,因此结果可能有多个,每一个都保存
    cursor = c.execute("SELECT title,entryPath FROM note WHERE title LIKE ?", (title,))
    input_file = []
    for row in cursor:
        title, path = row[0], row[1]
        # 没有用新版有道云笔记打开的脑图,本地文件不存在的情况忽略
        if path != None and os.path.exists(path):
            input_file.append((title, path))
    conn.close()
    return input_file

if __name__ == '__main__':
    if len(sys.argv) != 4:
        print(f'usage: {sys.argv[0]} db_path title save.xmind')
        exit(0)
    # 从有道云笔记的sqlite中得到文件路径
    input_file = get_mindmap_path(sys.argv[1], sys.argv[2])
    # 将这些有道云笔记脑图转换为xmind
    convert_mindmap_xmind(input_file, sys.argv[3])
