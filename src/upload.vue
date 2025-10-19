<script setup>
import Nav from './components/nav.vue'
import { onMounted } from 'vue';

onMounted(() => $('#uploadForm').on('submit', function (e) {
  e.preventDefault();

  const fileInput = $('#fileInput')[0];
  const file = fileInput.files[0];

  if (!file) {
    showMessage('请选择一个文件', 'danger');
    return;
  }

  const formData = new FormData();
  formData.append('file', file);

  // 显示进度条
  $('#progressContainer').show();
  $('#progressBar').css('width', '0%').text('0%');

  $.ajax({
    url: '/upload',
    type: 'POST',
    data: formData,
    processData: false,       // 不处理数据（交给 FormData）
    contentType: false,       // 不设置 Content-Type（浏览器自动设置）
    xhr: function () {
      // 自定义 XMLHttpRequest 以监听上传进度
      const xhr = new window.XMLHttpRequest();
      xhr.upload.addEventListener('progress', function (e) {
        if (e.lengthComputable) {
          const percent = Math.round((e.loaded / e.total) * 100);
          $('#progressBar')
            .css('width', percent + '%')
            .text(percent + '%');
        }
      });
      return xhr;
    },
    success: function (response) {
      showMessage('上传成功！', 'success');
      console.log('服务器响应:', response);
    },
    error: function (xhr, status, error) {
      console.log('服务器响应:', error);
      showMessage('上传失败: ' + (xhr.responseJSON?.message || error), 'danger');
    }
  });
}))

function showMessage(text, type) {
  $('#message')
    .removeClass('alert-success alert-danger')
    .addClass(`alert alert-${type}`)
    .text(text)
    .show();
}
</script>

<template>
  <header>
    <Nav></Nav>
  </header>
  <main class="container d-flex justify-content-center align-items-center" style="background-color: #f8f9fa;">
    <div>
      <h2>上传文件</h2>

      <form id="uploadForm" class="mt-4" method="post">
        <div class="mb-3">
          <label for="fileInput" class="form-label">选择文件</label>
          <input class="form-control" type="file" id="fileInput" name="file" required>
        </div>
        <button type="submit" class="btn btn-primary">上传</button>
      </form>

      <!-- 进度条容器（初始隐藏） -->
      <div id="progressContainer" class="mt-3" style="display:none;">
        <div class="progress" style="height: 25px;">
          <div id="progressBar" class="progress-bar" role="progressbar" style="width: 0%;">0%</div>
        </div>
      </div>

      <div id="message" class="mt-3"></div>
    </div>
  </main>

</template>