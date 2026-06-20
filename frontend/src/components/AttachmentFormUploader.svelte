<script>
  import { createEventDispatcher } from "svelte";
  
  let { files = [], disabled = false, maxFiles = 5 } = $props();
  const dispatch = createEventDispatcher();
  
  let dragActive = $state(false);
  let fileInput;
  
  // File validation
  const allowedTypes = [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'text/plain',
    'image/jpeg',
    'image/png',
    'image/gif',
    'application/zip'
  ];
  
  const allowedExtensions = [
    '.pdf', '.doc', '.docx', '.xls', '.xlsx',
    '.ppt', '.pptx', '.txt', '.jpg', '.jpeg',
    '.png', '.gif', '.zip'
  ];
  
  const maxFileSize = 10 * 1024 * 1024; // 10MB
  
  function validateFile(file) {
    // Check file size
    if (file.size > maxFileSize) {
      return { valid: false, error: `فایل ${file.name} بیش از 10 مگابایت است` };
    }
    
    // Check file type
    if (!allowedTypes.includes(file.type)) {
      return { valid: false, error: `نوع فایل ${file.name} مجاز نیست` };
    }
    
    // Check file extension
    const extension = '.' + file.name.split('.').pop().toLowerCase();
    if (!allowedExtensions.includes(extension)) {
      return { valid: false, error: `پسوند فایل ${file.name} مجاز نیست` };
    }
    
    return { valid: true };
  }
  
  function handleFiles(newFiles) {
    if (disabled) return;
    
    const fileArray = Array.from(newFiles);
    const validFiles = [];
    const errors = [];
    
    for (const file of fileArray) {
      // Check if we've reached max files
      if (files.length + validFiles.length >= maxFiles) {
        errors.push(`حداکثر ${maxFiles} فایل مجاز است`);
        break;
      }
      
      // Check if file already exists
      if (files.some(f => f.name === file.name && f.size === file.size)) {
        errors.push(`فایل ${file.name} قبلاً اضافه شده است`);
        continue;
      }
      
      const validation = validateFile(file);
      if (validation.valid) {
        validFiles.push(file);
      } else {
        errors.push(validation.error);
      }
    }
    
    if (validFiles.length > 0) {
      dispatch('filesAdded', { files: validFiles });
    }
    
    if (errors.length > 0) {
      dispatch('error', { message: errors.join('\n') });
    }
  }
  
  function handleDrop(e) {
    e.preventDefault();
    dragActive = false;
    
    const droppedFiles = e.dataTransfer.files;
    handleFiles(droppedFiles);
  }
  
  function handleDragOver(e) {
    e.preventDefault();
    if (!disabled) {
      dragActive = true;
    }
  }
  
  function handleDragLeave(e) {
    e.preventDefault();
    dragActive = false;
  }
  
  function handleFileSelect(e) {
    const selectedFiles = e.target.files;
    handleFiles(selectedFiles);
    
    // Reset input
    e.target.value = '';
  }
  
  function removeFile(index) {
    dispatch('fileRemoved', { index });
  }
  
  function formatFileSize(bytes) {
    if (bytes === 0) return '0 بایت';
    const k = 1024;
    const sizes = ['بایت', 'کیلوبایت', 'مگابایت'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }
  
  function getFileIcon(fileName) {
    const extension = fileName.split('.').pop().toLowerCase();
    
    switch (extension) {
      case 'pdf':
        return '📄';
      case 'doc':
      case 'docx':
        return '📝';
      case 'xls':
      case 'xlsx':
        return '📊';
      case 'ppt':
      case 'pptx':
        return '📋';
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif':
        return '🖼️';
      case 'zip':
        return '🗜️';
      case 'txt':
        return '📄';
      default:
        return '📎';
    }
  }
</script>

<div class="space-y-4">
  <!-- Upload Area -->
  <div
    class="border-2 border-dashed rounded-lg p-6 text-center transition-colors
    {dragActive ? 'border-brand-400 bg-brand-50' : 'border-slate-300'}
    {disabled ? 'opacity-50 cursor-not-allowed' : 'hover:border-brand-400 hover:bg-slate-50 cursor-pointer'}"
    ondrop={handleDrop}
    ondragover={handleDragOver}
    ondragleave={handleDragLeave}
    onclick={() => !disabled && fileInput?.click()}
  >
    <div class="space-y-2">
      <div class="text-4xl">📎</div>
      <div class="text-sm text-slate-600">
        {#if disabled}
          آپلود فایل غیرفعال است
        {:else}
          فایل‌ها را اینجا بکشید یا کلیک کنید
        {/if}
      </div>
      <div class="text-xs text-slate-500">
        حداکثر {maxFiles} فایل، هر فایل حداکثر 10 مگابایت
      </div>
      <div class="text-xs text-slate-400">
        فرمت‌های مجاز: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, TXT, JPG, PNG, GIF, ZIP
      </div>
    </div>
  </div>
  
  <!-- Hidden File Input -->
  <input
    bind:this={fileInput}
    type="file"
    multiple
    accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.jpg,.jpeg,.png,.gif,.zip"
    onchange={handleFileSelect}
    class="hidden"
    {disabled}
  />
  
  <!-- File List -->
  {#if files.length > 0}
    <div class="space-y-2">
      <h4 class="text-sm font-medium text-slate-700">فایل‌های انتخاب شده ({files.length})</h4>
      <div class="space-y-2">
        {#each files as file, index}
          <div class="flex items-center justify-between p-3 bg-slate-50 rounded-lg">
            <div class="flex items-center gap-3 flex-1 min-w-0">
              <span class="text-lg">{getFileIcon(file.name)}</span>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium text-slate-900 truncate">
                  {file.name}
                </div>
                <div class="text-xs text-slate-500">
                  {formatFileSize(file.size)}
                </div>
              </div>
            </div>
            {#if !disabled}
              <button
                onclick={() => removeFile(index)}
                class="p-1 text-slate-400 hover:text-rose-600 transition-colors"
                title="حذف فایل"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>