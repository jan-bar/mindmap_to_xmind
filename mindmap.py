import json
import xmind
import sys
import os


def convert_mindmap_xmind(input_file, save_file):
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

    if os.path.exists(save_file):
        os.remove(save_file)  # 目标文件存在则删除

    mind = xmind.load(save_file)
    sheet1 = mind.getPrimarySheet()
    # 获取sheet并设置名称为输入文件名
    sheet1.setTitle(os.path.basename(input_file))
    root1 = sheet1.getRootTopic()
    # 设置属性为: 逻辑图(向右), 和有道云笔记保持一致
    root1.setAttribute('structure-class', 'org.xmind.ui.logic.right')

    def set_topic(root, val_data):
        root.setTitle(val_data['topic'])
        for sub_val in val_data['sub_topic'].values():
            set_topic(root.addSubTopic(), sub_val)
    set_topic(root1, data_root)  # 递归设置子节点

    xmind.save(mind, save_file)


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print(f'usage: {sys.argv[0]} input.mindmap save.xmind')
        exit(0)
    convert_mindmap_xmind(sys.argv[1], sys.argv[2])
