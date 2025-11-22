<script setup>
import { ref } from 'vue';
import File from './file.vue';
import Dir from './dir.vue';
let filelist = ref([]);
let length = ref(0);
let dirlist = ref();
import { onMounted } from 'vue';

class RootPath {
    path;
    file;
    child;
    constructor(path) {
        this.path = path;
        this.file = [];
        this.child = new Map();
    }
}

onMounted(() =>
    $.ajax({
        url: '/list',
        type: 'POST',
        dataType: 'json',
        success: function (response) {
            console.log('成功:', response);
            let file = [];
            let rm = new Map();
            response.forEach(function (value, _) {
                if (!value.Path.includes("/")) { // 如果不是目录，就是文件
                    file.push(value);
                    length.value++;
                    return
                }

                // 处理在目录中的文件
                /* 构造目录树，结构类似
                                        a  （type=RootPath path=a and not contain file）
                                    /       \
                                    b       c   （type=RootPath path=a/b or a/c and contain file）
                                /      \  /     \
                                c      d  e      f
                 */
                let paths = value.Path.split("/");
                if (rm.get(paths[0]) == undefined) {
                    rm.set(paths[0], new RootPath(paths[0]))
                }
                let root = rm.get(paths[0]);
                let parent = root;
                let index = paths.length;
                for (let i = 1; i <= index; i++) {//遍历 root -> path的每一文件夹和它本身
                    if (i == index-1){ // 如果是文件
                        parent.file.push(value);
                        length.value++;
                        break
                    }
                    if (parent.child.get(paths[i]) == undefined) {
                        parent.child.set(paths[i], new RootPath(paths.slice(0,i+1).join('/')));
                    }
                    parent = parent.child.get(paths[i]);
                }
            });
            filelist.value = file;
            dirlist.value = rm;
            console.table(rm);
        },
        error: function (xhr, status, error) {
            console.error('错误:', error);
        }
    })
)
</script>

<template>
    <h2 class="text-center">所有文件（总数：{{ length }}）</h2>
    <div v-for="(item, index) in filelist" :key="item.Name">
        <File :path="item.Path" :size="item.Size" :updateAt="item.UpdatedAt"></File>
    </div>
    <div v-for="[index, item] in dirlist" :key="item.path">
        <Dir :path="item.path" :file="item.file" :child="item.child"></Dir>
    </div>
</template>
