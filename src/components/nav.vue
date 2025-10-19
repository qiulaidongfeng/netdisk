<script setup>
import {ref} from 'vue';
var Name = ref("");
var ID = ref("");

let id = Cookies.get("ID");
let b = typeof id != "undefined";
if (b){
  ID.value = id;
}

let name = Cookies.get("Name");
b = typeof name != "undefined";
if (b){
  Name.value = name;
}
let login = ref(Name.value != "");
</script>

<template>
    <nav class="navbar d-flex flex-column justify-content-between align-items-stretch" style="height: 100vh;background-color:  #f8f9fa; max-width: 20vw;">
        <!-- 顶部：折叠按钮 -->
        <button class="navbar-toggler invisible-lg-but-space align-self-start" type="button" data-bs-toggle="collapse"
            data-bs-target="#verticalNavbar">
            <span class="navbar-toggler-icon"></span>
        </button>

        <!-- 底部：所有导航内容（Logo + 菜单） -->
        <div class="collapse show" id="verticalNavbar">
            <p class="navbar-brand mb-5">简易网盘</p>
            <ul class="navbar-nav flex-column">
                <li class="nav-item"><a class="nav-link" href="/">首页</a></li>
                <li class="nav-item"><a class="nav-link" href="/register.html">注册</a></li>
                <li class="nav-item"><a class="nav-link" href="/login.html">登录</a></li>
                <li class="nav-item"><a class="nav-link" href="/user.html">账号管理</a></li>
                <li class="nav-item"><a class="btn btn-primary" href="/upload.html">上传文件</a></li>
            </ul>
            <p class="text-break" v-if="login">用户名：{{ name }}</p>
            <p class="text-break" v-else>未登录</p>
            <p class="text-break" v-if="login">用户Id：{{ id }}</p>
            <p class="text-break" v-else>未登录</p>
        </div>
    </nav>
</template>

<style>
.invisible-lg-but-space {
    visibility: visible !important;
}

@media (min-width: 992px) {
    .invisible-lg-but-space {
        visibility: hidden !important;
    }
}
</style>