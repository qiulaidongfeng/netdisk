<script setup>
defineProps({
  "info": {
    type: String,
    required: true
  },
  "id": {
    type: String,
    required: true
  }
})
const togglePassword = (isShow, id) => {
  const input = document.getElementById(id);
  if (input) {
    input.type = isShow ? 'text' : 'password';
  }
};
</script>

<template>
  <div class="d-flex flex-phone-column d-lg-block">
    <span>{{ info }}</span><input type="password" class="form-control d-inline-block w-auto" :name="id" :id="id" required>
    <div class="password mt-2 d-inline-block">
      <span @click="togglePassword(true, id)">显示</span>
      <span @click="togglePassword(false, id)">隐藏</span>
    </div>
  </div>
</template>

<style scoped>
.password span {
  display: inline-block;
  border: 1px solid black;
  margin-left: 4px;
  text-align: center;
  width: 4rem;
  color: black;
}

.password span:hover {
  transform: translateY(-5px);
}
</style>