<script setup>
import Nav from './components/nav.vue'
import Password from './components/password.vue'
</script>

<template>
  <header>
    <Nav></Nav>
  </header>
  <main class="container d-flex justify-content-center align-items-center" style="background-color: #f8f9fa;">
    <div class="register d-flex p-3 flex-column border border-1">

      <div class="d-flex mb-3 justify-content-center align-items-center" style="background-color: #d0e7ff;">
        <div class="register-banner">
          <p>登录</p>
        </div>
      </div>
      <form action="/login" method="post" class="mb-3">
        <div class="mb-3">
          <span>用户ID：</span><input type="text" class="form-control d-inline-block w-auto" name="id" id="id" required>
        </div>
        <div class="mb-3">
          <Password info="密码：" id="password1"></Password>
          <br>
          <button type="submit" class="btn btn-primary"> 登录 </button>
        </div>
      </form>
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