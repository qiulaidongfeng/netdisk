<script setup>
defineProps({
    "path":{
        type:String,
        required:true
    },
    "size":{
        type:Number,
        required:true
    },
    "updateAt":{
        type:String,
        required:true
    }
})

function fmt(v) {
    if ((v / 1024 / 1024) < 1) {
        return (v / 1024).toFixed(2) + "kb"
    }
    return (v / 1024 / 1024).toFixed(2) + "Mb"
}

</script>

<template>
    <div class="d-flex justify-content-between">
        <p>{{ path }}</p>
        <p>{{ fmt(size) }}</p>
    </div>
    <div class="d-flex justify-content-between flex-wrap">
        <a :href="`/download/${path}/`" class="btn btn-primary">下载</a>
        <a :href="`/delete/${path}/`" class="btn btn-primary">删除</a>
        <p>修改日期：{{ updateAt }}</p>
    </div>
    <hr>
</template>

<style>
p {
    margin-bottom: 0.5rem;
    word-break: break-all;
}
</style>