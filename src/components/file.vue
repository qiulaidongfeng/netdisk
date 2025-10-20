<script setup>
import { ref } from 'vue';
let list = ref();
import { onMounted } from 'vue';

function fmt(v) {
  if ((v / 1024 / 1024) < 1) {
    return (v / 1024).toFixed(2) + "kb"
  }
  return (v / 1024 / 1024).toFixed(2) + "Mb"
}

onMounted(() =>
    $.ajax({
        url: '/list',
        type: 'POST',
        dataType: 'json',
        success: function (response) {
            console.log('成功:', response);
            list.value = response;
        },
        error: function (xhr, status, error) {
            console.error('错误:', error);
        }
    })
)
</script>

<template>
    <h2 class="text-center">所有文件（总数：{{ list && list.length || 0 }}）</h2>
    <div v-for="(item, index) in list" :key="item.Path">
        <div class="d-flex justify-content-between">
            <p>{{ item.Path }}</p>
            <p>{{ fmt(item.Size) }}</p>
        </div>
        <div class="d-flex justify-content-between flex-wrap">
            <a :href="`/download/${item.Path}/`" class="btn btn-primary">下载</a>
            <a :href="`/delete/${item.Path}/`" class="btn btn-primary">删除</a>
            <br class="d-block d-lg-none">
            <p style="text-align: right;">修改日期：{{ item.UpdatedAt }}</p>
        </div>
    </div>
</template>

<style>
p {
    margin-bottom: 0.5rem;
}
</style>