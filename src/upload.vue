<script setup>
import Nav from './components/nav.vue'
import { onMounted, ref } from 'vue';

const isUploading = ref(false);
const totalFiles = ref(0);
const uploadedFiles = ref(0);

onMounted(() => {

  // --- 文件选择处理 (Files Input) ---
  $('#fileInputFiles').on('change', function(e) {
      console.log("文件选择已更改 (Files):", e.target.files.length, "个文件");
  });

  // --- 文件夹选择处理 (Folder Input) ---
  $('#fileInputFolder').on('change', function(e) {
      console.log("文件夹选择已更改 (Folder):", e.target.files.length, "个项目");
      // 对于 webkitdirectory，webkitRelativePath 即为所需相对路径，无需额外计算。
  });

  // --- 文件上传表单提交 (文件) ---
  $('#uploadFormFiles').on('submit', async function (e) {
    e.preventDefault();

    if (isUploading.value) {
      showMessage('正在上传中，请稍候...', 'warning');
      return;
    }

    const fileInput = $('#fileInputFiles')[0];
    const files = Array.from(fileInput.files);

    if (files.length === 0) {
      showMessage('请选择至少一个文件', 'danger');
      return;
    }

    isUploading.value = true;
    totalFiles.value = files.length;
    uploadedFiles.value = 0;
    let hasErrorOccurred = false;

    showMessage('', 'info');
    showErrorMessage('', false);

    showMessage(`准备上传 ${totalFiles.value} 个文件...`, 'info');
    $('#progressContainer').show();
    updateProgressBar();

    for (const file of files) {
      if (hasErrorOccurred) {
          console.log(`因之前上传失败，已停止上传后续文件。`);
          break;
      }

      try {
        // 对于普通文件上传，使用文件本身的名称作为目标路径
        await uploadSingleFile(file, file.name);
        uploadedFiles.value++;
        updateProgressBar();
      } catch (error) {
        console.error(`上传文件 ${file.name} 失败:`, error);
        showErrorMessage(`文件 ${file.name} 上传失败: ${error.message || error}`, true);
        hasErrorOccurred = true;
      }
    }

    finalizeUpload(hasErrorOccurred);
    $('#fileInputFiles').val(''); // Clear input after upload
  });


  // --- 文件夹上传表单提交 (文件夹) ---
  $('#uploadFormFolder').on('submit', async function (e) {
    e.preventDefault();

    if (isUploading.value) {
      showMessage('正在上传中，请稍候...', 'warning');
      return;
    }

    const fileInput = $('#fileInputFolder')[0];
    const files = Array.from(fileInput.files);

    if (files.length === 0) {
      showMessage('请选择一个包含文件的文件夹', 'danger');
      return;
    }

    isUploading.value = true;
    totalFiles.value = files.length;
    uploadedFiles.value = 0;
    let hasErrorOccurred = false;

    showMessage('', 'info');
    showErrorMessage('', false);

    showMessage(`准备上传文件夹中的 ${totalFiles.value} 个项目...`, 'info');
    $('#progressContainer').show();
    updateProgressBar();

    for (const file of files) {
      if (hasErrorOccurred) {
          console.log(`因之前上传失败，已停止上传后续文件。`);
          break;
      }

      try {
        // --- 关键逻辑 ---
        // 对于文件夹上传，直接使用浏览器提供的 webkitRelativePath
        let targetServerPath = file.name; // Default fallback
        if (file.webkitRelativePath) {
            targetServerPath = file.webkitRelativePath;
            console.log(`文件夹上传: 文件 ${file.name} 将使用相对路径 "${targetServerPath}" 上传`);
        } else {
             console.warn(`文件夹上传: 文件 ${file.name} 缺少 webkitRelativePath，将使用原始文件名`);
        }

        // --- 将计算出的相对路径作为查询参数传递给后端 ---
        await uploadSingleFile(file, targetServerPath);
        uploadedFiles.value++;
        updateProgressBar();
      } catch (error) {
        console.error(`上传文件 ${file.name} 失败:`, error);
        showErrorMessage(`文件 ${file.name} 上传失败: ${error.message || error}`, true);
        hasErrorOccurred = true;
      }
    }

    finalizeUpload(hasErrorOccurred);
    $('#fileInputFolder').val(''); // Clear input after upload
  });


  /**
   * 更新全局上传进度条
   */
  function updateProgressBar() {
    if (totalFiles.value > 0) {
      const percent = Math.round((uploadedFiles.value / totalFiles.value) * 100);
      $('#progressBar')
        .css('width', percent + '%')
        .attr('aria-valuenow', percent)
        .text(`${uploadedFiles.value}/${totalFiles.value} (${percent}%)`);
    } else {
      $('#progressBar').css('width', '0%').attr('aria-valuenow', 0).text('0%');
    }
  }

  /**
   * 上传单个文件到服务器
   * @param {File} file - 要上传的文件对象
   * @param {string} targetServerPath - 在服务器上使用的文件路径（可以包含路径分隔符 '/'）
   *                                   对于普通文件上传，通常是 file.name
   *                                   对于文件夹上传，是 file.webkitRelativePath
   * @returns {Promise} 上传成功则 resolve，失败则 reject
   *
   * 修改说明：
   * 此函数将 targetServerPath 作为查询参数 '?path=' 附加到上传 URL (/upload) 上。
   * 文件内容通过 FormData 发送，但不依赖 filename 参数传递路径。
   * 后端的 /upload 接口需要从 ctx.Query("path") 读取此路径。
   */
  function uploadSingleFile(file, targetServerPath) {
    return new Promise((resolve, reject) => {
      const formData = new FormData();

      // 1. 添加文件内容到 FormData
      //    不再需要特别指定 filename，让浏览器使用默认值即可
      formData.append('file', file); // Standard append without custom filename

      // 2. 构造带查询参数的完整上传 URL
      //    使用 encodeURIComponent 对路径进行编码，以防包含特殊字符
      const encodedPath = encodeURIComponent(targetServerPath);
      const uploadUrl = `/upload?path=${encodedPath}`;

      console.log(`准备上传文件 [${file.name}] 到 URL [${uploadUrl}]`);

      $.ajax({
        url: uploadUrl, // <-- 使用带 path 查询参数的 URL
        type: 'POST',
        data: formData,
        processData: false,
        contentType: false,
        success: function (response) {
          console.log(`文件 "${file.name}" (路径: '${targetServerPath}') 上传成功:`, response);
          resolve(response);
        },
        error: function (xhr, status, error) {
          console.error(`文件 "${file.name}" (路径: '${targetServerPath}') 上传失败:`, status, error);
          const errorMsg = xhr.responseJSON?.message || xhr.statusText || error.toString();
          reject(new Error(errorMsg));
        }
      });
    });
  }

  /**
   * 通用的上传结束处理逻辑
   * @param {boolean} hasErrorOccurred - 是否发生了错误
   */
  function finalizeUpload(hasErrorOccurred) {
    isUploading.value = false;
    $('#progressContainer').hide();

    if (!hasErrorOccurred && uploadedFiles.value === totalFiles.value) {
      showMessage(`所有 ${totalFiles.value} 个项目上传成功！`, 'success');
    } else if (hasErrorOccurred) {
      showMessage(`上传因错误中断。已完成 ${uploadedFiles.value}/${totalFiles.value} 个。请查看下方详细错误信息。`, 'warning');
    } else {
      showMessage(`上传结束，但状态不明。`, 'info');
    }
  }


  /**
   * 显示通用消息
   * @param {string} text - 消息文本
   * @param {'success'|'danger'|'info'|'warning'} type - 消息类型
   */
  function showMessage(text, type) {
    const $msgElement = $('#message');
    $msgElement
      .removeClass('alert-success alert-danger alert-info alert-warning')
      .addClass(`alert alert-${type}`)
      .text(text);

     if (text) {
         $msgElement.show();
         // Auto-hide success/info messages after some time
         if (type === 'success' || type === 'info') {
            setTimeout(() => {
              $msgElement.fadeOut();
            }, 5000);
         }
     } else {
         $msgElement.hide();
     }
  }

  /**
   * 显示错误消息
   * @param {string} text - 错误消息文本
   * @param {boolean} append - 是否追加到现有错误信息后面
   */
  function showErrorMessage(text, append = false) {
    const $errorMsgElement = $('#errorMessage');
    if (append && text) {
        const currentText = $errorMsgElement.text();
        const newText = currentText ? `${currentText}\n${text}` : text;
        $errorMsgElement.text(newText);
        $errorMsgElement.show();
    } else if (!append) {
        $errorMsgElement.text(text);
        if (text) {
            $errorMsgElement.show();
        } else {
            $errorMsgElement.hide();
        }
    }
  }

});
</script>

<template>
  <header>
    <Nav></Nav>
  </header>
  <main class="container d-flex flex-column align-items-center py-4" style="min-height: 80vh; background-color: #f8f9fa;">
    <div class="w-100" style="max-width: 600px;">
      <h2 class="mb-4 text-center">上传文件/文件夹</h2>

      <!-- 文件上传区域 -->
      <form id="uploadFormFiles" class="mt-4 mb-5 border p-4 rounded bg-white shadow-sm">
        <h5 class="mb-3 text-center">上传文件</h5>
        <div class="mb-3">
          <label for="fileInputFiles" class="form-label">选择文件</label>
          <input class="form-control" type="file" id="fileInputFiles" name="file" multiple>
          <div class="form-text">可以选择单个或多个文件。</div>
        </div>
        <div class="d-flex justify-content-center">
           <button type="submit" class="btn btn-primary" :disabled="isUploading">
             {{ isUploading ? '上传中...' : '开始上传文件' }}
           </button>
        </div>
      </form>

      <!-- 文件夹上传区域 -->
      <form id="uploadFormFolder" class="mt-4 border p-4 rounded bg-white shadow-sm">
       <h5 class="mb-3 text-center">上传文件夹</h5>
        <div class="mb-3">
          <label for="fileInputFolder" class="form-label">选择文件夹</label>
          <input class="form-control" type="file" id="fileInputFolder" name="file" webkitdirectory>
          <div class="form-text">选择一个文件夹，将上传其内部所有文件及子文件夹（保持原有目录结构）。</div>
        </div>
        <div class="d-flex justify-content-center">
           <button type="submit" class="btn btn-primary" :disabled="isUploading">
             {{ isUploading ? '上传中...' : '开始上传文件夹' }}
           </button>
        </div>
      </form>

      <!-- 进度条容器（默认隐藏） -->
      <div id="progressContainer" class="mt-4" style="display:none;">
        <div class="progress" style="height: 25px;">
          <div id="progressBar" class="progress-bar" role="progressbar" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100" style="width: 0%;">0%</div>
        </div>
      </div>

      <!-- 通用消息显示区域 -->
      <div id="message" class="mt-4"></div> <!-- Added ID for JS -->

      <!-- 独立的错误消息显示区域 -->
      <div id="errorMessage" class="mt-3 alert alert-danger" style="white-space: pre-line; display: none;"></div> <!-- Added ID for JS -->

    </div>
  </main>
</template>

<style scoped>
/* Ensure the buttons are centered */
.d-flex.justify-content-center {
  margin-top: 1rem; /* Add some space above the button */
}

/* Optional: Style the info message differently if needed */
.alert-info {
  color: #0c5460;
  background-color: #d1ecf1;
  border-color: #bee5eb;
}

/* Optional: Style the warning message */
.alert-warning {
  color: #856404;
  background-color: #fff3cd;
  border-color: #ffeaa7;
}
</style>