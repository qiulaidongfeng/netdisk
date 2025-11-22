<script setup>
import File from './file.vue';
defineOptions({
    name: 'Dir'
});
const props = defineProps({
    "path": {
        type: String,
        required: true
    },
    "file": {
        type: Array,
        required: true
    },
    "child": {
        type: Map,
        required: true
    },
})
let index = props.path; 
</script>

<template>
    <div>
        <div class="d-flex justify-content-between">
            <p>文件夹:{{ path }}</p>
            <button class="mt-2 btn btn-primary" type="button" data-bs-toggle="collapse" :data-bs-target="`#${index}`"
                aria-expanded="false" aria-controls="`#${index}`">
                展开
            </button>
        </div>
        <div class="collapse" :id="`${index}`">
            <div v-for="(item, index) in file" :key="item.Path">
                <File :path="item.Path" :size="item.Size" :updateAt="item.UpdatedAt"></File>
            </div>
            <div v-for="[index, item] in child" :key="item.path">
                <Dir :path="item.path" :index="item.path" :file="item.file" :child="item.child"></Dir>
            </div>
        </div>
    </div>
</template>

<style>
hr {
    margin-bottom: 0.5rem;
    margin-top: 0.5rem;
}
</style>
