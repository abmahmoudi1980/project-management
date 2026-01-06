<script>
  import { createEventDispatcher } from "svelte";
  import { api } from "../lib/api.js";

  let { taskId, disabled = false } = $props();
  const dispatch = createEventDispatcher();

  let isDragOver = $state(false);
  let isUploading = $state(false);
  let uploadProgress = $state({});
  let validationErrors = $state([]);
  let fileInput;

  // File validation configuration
  const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
  const MAX_TOTAL_SIZE = 100 * 1024 * 1024; // 100MB
  const ALLOWED_EXTENSIONS = ['.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx', '.txt', '.jpg', '.jpeg', '.png', '.gif', '.zip'];
  const BLOCKED_EXTENSIONS = ['.exe', '.bat', '.cmd', '.scr', '.pif', '.com', '.js', '.vbs', '.jar'];

  function validateFile(file) {
    const errors = [];
    const extension = '.' + file.name.split('.').pop().toLowerCase();
    
    // Check blocked extensions
    if (BLOCKED_EXTENSIONS.includes(extension)) {
      errors.push(`فایل ${file.name}: نوع فایل مجاز نیست`);
      return errors;
    }
    
    // Check allowed extensions
    if (!ALLOWED_EXTENSIONS.includes(extension)) {
      errors.push(`فایل ${file.name}: نوع فایل پشتیبانی نمی‌شود`);
    }
    
    // Check file size
    if (file.size > MAX_FILE_SIZE) {
      const sizeMB = (MAX_FILE_SIZE / (1024 * 1024)).toFixed(0);
      errors.push(`فایل ${file.name}: حجم فایل نباید بیشتر از ${sizeMB} مگابایت باشد`);
    }
    
    return errors;
  }

  function validateFiles(files) {
    const errors = [];
    let totalSize = 0;
    
    for (const file of files) {
      const fileErrors = validateFile(file);
      errors.push(...fileErrors);
      totalSize += file.size;
    }
    
    // Check total size
    if (totalSize > MAX_TOTAL_SIZE) {
      const sizeMB = (MAX_TOTAL_SIZE / (1024 * 1024)).toFixed(0);
      errors.push(`مجموع حجم فایل‌ها نباید بیشتر از ${sizeMB} مگابایت باشد`);
    }
    
    return errors;
  }

  function formatFileSize(bytes) {
    if (bytes === 0) return '0 بایت';
    const k = 1024;
    const sizes = ['بایت', 'کیلوبایت', 'مگابایت', 'گیگابایت'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  async function uploadFiles(files) {
    if (!files || files.length === 0) return;
    
    // Validate files
    const errors = validateFiles(files);
    if (errors.length > 0) {
      validationErrors = errors;
      return;
    }
    
    validationErrors = [];
    isUploading = true;
    uploadProgress = {};
    
    try {
      // Initialize progress for all files
      for (const file of files) {
        uploadProgress[file.name] = { progress: 0, status: 'uploading' };
      }
      uploadProgress = { ...uploadProgress };
      
      // Use the centralized API client
      const response = await api.attachments.upload(taskId, files);
      
      // Update progress for successful uploads
      for (const file of files) {
        uploadProgress[file.name] = { progress: 100, status: 'completed' };
      }
      uploadProgress = { ...uploadProgress };
      
      // Clear progress after a delay
      setTimeout(() => {
        uploadProgress = {};
        isUploading = false;
      }, 2000);
      
      dispatch('uploaded', response);
      
    } catch (error) {
      console.error('Upload error:', error);
      
      // Update progress for failed uploads
      for (const file of files) {
        uploadProgress[file.name] = { progress: 0, status: 'error' };
      }
      uploadProgress = { ...uploadProgress };
      
      dispatch('error', { message: error.message || 'خطا در آپلود فایل' });
      isUploading = false;
      
      // Clear progress after a delay
      setTimeout(() => {
        uploadProgress = {};
      }, 3000);
    }
  }

  function handleDragOver(e) {
    e.preventDefault();
    if (!disabled && !isUploading) {
      isDragOver = true;
    }
  }

  function handleDragLeave(e) {
    e.preventDefault();
    isDragOver = false;
  }

  function handleDrop(e) {
    e.preventDefault();
    isDragOver = false;
    
    if (disabled || isUploading) return;
    
    const files = Array.from(e.dataTransfer.files);
    uploadFiles(files);
  }

  function handleFileSelect(e) {
    if (disabled || isUploading) return;
    
    const files = Array.from(e.target.files);
    uploadFiles(files);
    
    // Clear input for next selection
    e.target.value = '';
  }

  function cancelUpload() {
    // In a real implementation, you would cancel the XMLHttpRequest
    isUploading = false;
    uploadProgress = {};
    validationErrors = [];
  }

  function openFileDialog() {
    if (!disabled && !isUploading) {
      fileInput.click();
    }
  }
</script>

<div class="space-y-4">
  <!-- Upload Area -->
  <div
    class="border-2 border-dashed rounded-lg p-6 text-center transition-colors
    {isDragOver ? 'border-indigo-500 bg-indigo-50' : 'border-slate-300'}
    {disabled || isUploading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer hover:border-indigo-400 hover:bg-slate-50'}"
    ondragover={handleDragOver}
    ondragleave={handleDragLeave}
    ondrop={handleDrop}
    onclick={openFileDialog}
  >
    <input
      bind:this={fileInput}
      type="file"
      multiple
      accept={ALLOWED_EXTENSIONS.join(',')}
      onchange={handleFileSelect}
      class="hidden"
      {disabled}
    />
    
    <div class="space-y-3">
      <div class="mx-auto w-12 h-12 text-slate-400">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
      </div>
      
      <div>
        <p class="text-lg font-medium text-slate-700">
          {isUploading ? 'در حال آپلود...' : 'فایل‌ها را اینجا بکشید یا کلیک کنید'}
        </p>
        <p class="text-sm text-slate-500 mt-1">
          حداکثر {formatFileSize(MAX_FILE_SIZE)} برای هر فایل، مجموعاً {formatFileSize(MAX_TOTAL_SIZE)}
        </p>
        <p class="text-xs text-slate-400 mt-1">
          فرمت‌های مجاز: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, TXT, JPG, PNG, GIF, ZIP
        </p>
      </div>
    </div>
  </div>

  <!-- Validation Errors -->
  {#if validationErrors.length > 0}
    <div class="bg-rose-50 border border-rose-200 rounded-lg p-4">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="w-5 h-5 text-rose-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="mr-3">
          <h3 class="text-sm font-medium text-rose-800">خطاهای اعتبارسنجی</h3>
          <div class="mt-2 text-sm text-rose-700">
            <ul class="list-disc list-inside space-y-1">
              {#each validationErrors as error}
                <li>{error}</li>
              {/each}
            </ul>
          </div>
        </div>
      </div>
    </div>
  {/if}

  <!-- Upload Progress -->
  {#if Object.keys(uploadProgress).length > 0}
    <div class="bg-slate-50 border border-slate-200 rounded-lg p-4">
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-sm font-medium text-slate-700">پیشرفت آپلود</h3>
        {#if isUploading}
          <button
            onclick={cancelUpload}
            class="text-xs text-slate-500 hover:text-slate-700 underline"
          >
            لغو
          </button>
        {/if}
      </div>
      
      <div class="space-y-3">
        {#each Object.entries(uploadProgress) as [fileName, progress]}
          <div class="space-y-1">
            <div class="flex items-center justify-between text-sm">
              <span class="text-slate-700 truncate flex-1 ml-2">{fileName}</span>
              <span class="text-slate-500 text-xs">
                {progress.status === 'uploading' ? `${Math.round(progress.progress)}%` : 
                 progress.status === 'processing' ? 'در حال پردازش...' :
                 progress.status === 'completed' ? 'تکمیل شد' : 'خطا'}
              </span>
            </div>
            <div class="w-full bg-slate-200 rounded-full h-2">
              <div
                class="h-2 rounded-full transition-all duration-300
                {progress.status === 'completed' ? 'bg-emerald-500' : 
                 progress.status === 'error' ? 'bg-rose-500' : 'bg-indigo-500'}"
                style="width: {progress.progress}%"
              ></div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>