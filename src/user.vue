<script setup>
import Nav from './components/nav.vue'
import Password from './components/password.vue'
import { onMounted, ref } from 'vue';

var Name = ref("");

let name = Cookies.get("Name");
let b = typeof name != "undefined";
if (b){
  Name.value = name;
}

let percent = ref("0%");
let limit = ref("50Mb");
let used = ref("0Mb");

function fmt(v) {
  if ((v / 1024 / 1024) < 1) {
    return (v / 1024).toFixed(2) + "kb"
  }
  return (v / 1024 / 1024).toFixed(2) + "Mb"
}

onMounted(() => {
  $.ajax({
    url: '/stat',
    type: 'POST',
    data: "json",
    success: function (response) {
      console.log(response);
      percent.value = (response.Used / response.Limit).toFixed(5) + "%";
      limit.value = fmt(response.Limit);
      used.value = fmt(response.Used);
    },
    error: function (xhr, status, error) {
      console.error(status);
      console.error(error);
    }
  }
  )
})

onMounted(()=>
  $("#password_form").submit(function (e) {
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
  <main class="container d-flex flex-column flex-wrap" style="background-color: #f8f9fa;">
    <h1 class="text-center mb-5">用户管理</h1>
    <form action="/setname" class="d-flex flex-wrap" method="post">
      <div class="d-flex d-lg-block flex-phone-column flex-wrap justify-content-between mb-3">
        <p>用户名：{{ Name }}</p>
        <div class="d-flex flex-wrap">
          <span>新用户名：</span><input class="form-control d-inline-block w-auto" type="text" name="name" required>
        </div>
        <button type="submit" class="btn btn-primary mt-2"> 修改 </button>
      </div>
    </form>

    <hr class="mb-4">

    <form id="password_form" class="d-flex flex-wrap" action="/set_password" method="post">
      <div class="d-flex d-lg-block flex-phone-column flex-wrap justify-content-between mb-3">
        <p class="p-1">密码</p>

        <Password info="新密码：" id="password1"></Password>

        <Password info="重复新密码：" id="password2"></Password>

        <button type="submit" class="btn btn-primary mt-2"> 修改 </button>
      </div>
    </form>

    <hr class="mb-4">
    <div class="d-flex flex-wrap justify-content-between mb-3">
      <span class="me-auto">已使用容量：{{ used }}</span>
      <span class="me-auto">总可用容量：{{ limit }}</span>
      <span class="me-auto">使用百分比：{{ percent }}</span>
    </div>
  </main>

</template>
