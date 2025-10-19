<script setup>
import Nav from './components/nav.vue'
import Password from './components/password.vue'
import { onMounted } from 'vue'

onMounted(()=>
  $("#form").submit(function (e) {
    e.preventDefault();
    if ($("#password1").val() != $("#password2").val()) {
      alert("两次输入的密码应该一致");
      return
    }
    this.submit();
  })
)
</script>

<template>
  <header>
    <Nav></Nav>
  </header>
  <main class="container d-flex justify-content-center align-items-center" style="background-color: #f8f9fa;">
    <div class="register d-flex p-3 flex-column border border-1">

      <div class="d-flex mb-3 justify-content-center align-items-center" style="background-color: #d0e7ff;">
        <div class="register-banner">
          <p>创建新账户</p>
          <p style="margin-left: -2rem;">填写以下信息完成注册</p>
        </div>
      </div>

      <form id="form" action="/register" method="post" class="mb-3">
        <div class="mb-3">
          <span>用户名：</span><input type="text" class="form-control d-inline-block w-auto" name="name" id="name" required>
        </div>
        <div class="mb-3">
          <Password info="密码：" id="password1"></Password>
          <br>
          <Password info="重复输入密码：" id="password2"></Password>
          <button type="submit" class="btn btn-primary"> 注册 </button>
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

/* 伪元素：向左扩展背景 */
.register-banner::before {
  content: '';
  position: absolute;
  top: 0;
  left: -2rem;
  /* 和负 margin 一致 */
  right: 0;
  /* 等价于 width: 100% + 2rem */
  bottom: 0;
  background-color: #d0e7ff;
  z-index: -1;
  /* 在内容下方 */
}
</style>