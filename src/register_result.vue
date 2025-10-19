<script setup>
import { onMounted,ref } from 'vue';
import Nav from './components/nav.vue'
var ID = ref("");

onMounted(() => {
  let id = Cookies.get("ID");
  console.log(id);
  let b = typeof id != "undefined";
  if (b) {
    ID.value = id;
  }
})
</script>

<template>
  <header>
    <Nav></Nav>
  </header>
  <main class="container d-flex justify-content-center align-items-center" style="background-color: #f8f9fa;">
    <div class="register d-flex p-3 flex-column border border-1">

      <div class="d-flex mb-3 justify-content-center align-items-center" style="background-color: #d0e7ff;">
        <div class="register-banner">
          <p>注册成功</p>
        </div>
      </div>
      <div class="mb-3">
        <span>用户ID：{{ ID }}</span>
      </div>
      <div class="mb-3">
        <a href="/" class="btn btn-primary"> 返回首页 </a>
      </div>
    </div>
  </main>

</template>

<style scoped>
.register {
  width: min(80vw, 600px);
}

.register-banner {
  background-color: #d0e7ff;
  position: relative;
  /* 确保伪元素能覆盖内容下方 */
  z-index: 0;
}
</style>