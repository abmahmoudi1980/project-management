<script>
  import { createEventDispatcher } from "svelte";
  import { api } from "../lib/api.js";
  import moment from "jalali-moment";

  let { attachments = [], currentUser = null, canDelete = false } = $props();
  const dispatch = createEventDispatcher();

  function formatFileSize(bytes) {
    if (bytes === 0) return '0 بایت';
    const k = 1024;
    const sizes = ['بایت', 'کیلوبایت', 'مگابایت', 'گیگابایت'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  function formatJalaliDate(dateString) {
    if (!dateString) return "";
    return moment(dateString).locale("fa").format("YYYY/MM/DD HH:mm");
  }

  function getFileIcon(filename, mimeType) {
    const extension = filename.split('.').pop().toLowerCase();
    
    // Image files
    if (['jpg', 'jpeg', 'png', 'gif'].includes(extension)) {
      return 'image';
    }
    
    // Document files
    if (['pdf'].includes(extension)) {
      return 'pdf';
    }
    
    if (['doc', 'docx'].includes(extension)) {
      return 'document';
    }
    
    if (['xls', 'xlsx'].includes(extension)) {
      return 'spreadsheet';
    }
    
    if (['ppt', 'pptx'].includes(extension)) {
      return 'presentation';
    }
    
    if (['zip', 'rar', '7z'].includes(extension)) {
      return 'archive';
    }
    
    if (['txt'].includes(extension)) {
      return 'text';
    }
    
    return 'file';
  }

  function getIconSvg(iconType) {
    const icons = {
      image: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />`,
      pdf: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />`,
      document: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />`,
      spreadsheet: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0V4a2 2 0 012-2h14a2 2 0 012 2v16a2 2 0 01-2 2H5a2 2 0 01-2-2z" />`,
      presentation: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m0 0V3a1 1 0 011 1v10a1 1 0 01-1 1H8a1 1 0 01-1-1V4a1 1 0 011-1h8z" />`,
      archive: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />`,
      text: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />`,
      file: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />`
    };
    return icons[iconType] || icons.file;
  }

  function handleDownload(attachment) {
    // Use the centralized API client to get download URL
    const downloadUrl = api.attachments.download(attachment.id);
    
    // Create temporary link and trigger download
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = attachment.original_filename;
    link.style.display = 'none';
    
    // Add credentials for authenticated download
    fetch(downloadUrl, {
      credentials: 'include'
    }).then(response => {
      if (response.ok) {
        return response.blob();
      }
      throw new Error('دانلود فایل با خطا مواجه شد');
    }).then(blob => {
      const url = window.URL.createObjectURL(blob);
      link.href = url;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    }).catch(error => {
      console.error('Download error:', error);
      dispatch('error', { message: error.message });
    });
  }

  function handleDelete(attachment) {
    if (confirm(`آیا از حذف فایل "${attachment.original_filename}" اطمینان دارید؟`)) {
      dispatch('delete', attachment);
    }
  }

  function getThumbnailUrl(attachment) {
    if (attachment.has_thumbnail) {
      return api.attachments.getThumbnail(attachment.id);
    }
    return null;
  }

  function isImage(filename) {
    const extension = filename.split('.').pop().toLowerCase();
    return ['jpg', 'jpeg', 'png', 'gif'].includes(extension);
  }
</script>

{#if attachments.length === 0}
  <div class="text-center py-8 text-slate-500">
    <div class="mx-auto w-12 h-12 text-slate-300 mb-3">
      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.586-6.586a2 2 0 00-2.828-2.828z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M9 12l2 2 4-4" />
      </svg>
    </div>
    <p>هیچ فایل پیوستی وجود ندارد</p>
  </div>
{:else}
  <div class="space-y-3">
    {#each attachments as attachment (attachment.id)}
      <div class="flex items-center gap-4 p-4 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors">
        <!-- File Icon/Thumbnail -->
        <div class="flex-shrink-0">
          {#if isImage(attachment.original_filename) && attachment.has_thumbnail}
            <img
              src={getThumbnailUrl(attachment)}
              alt={attachment.original_filename}
              class="w-12 h-12 object-cover rounded-lg border border-slate-200"
              onerror={(e) => {
                e.target.style.display = 'none';
                e.target.nextElementSibling.style.display = 'block';
              }}
            />
            <div class="w-12 h-12 bg-slate-200 rounded-lg flex items-center justify-center text-slate-400" style="display: none;">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                {@html getIconSvg(getFileIcon(attachment.original_filename, attachment.mime_type))}
              </svg>
            </div>
          {:else}
            <div class="w-12 h-12 bg-slate-200 rounded-lg flex items-center justify-center text-slate-400">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                {@html getIconSvg(getFileIcon(attachment.original_filename, attachment.mime_type))}
              </svg>
            </div>
          {/if}
        </div>

        <!-- File Info -->
        <div class="flex-1 min-w-0">
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0">
              <h4 class="text-sm font-medium text-slate-900 truncate">
                {attachment.original_filename}
              </h4>
              <div class="flex items-center gap-4 mt-1 text-xs text-slate-500">
                <span>{formatFileSize(attachment.file_size)}</span>
                <span>{formatJalaliDate(attachment.created_at)}</span>
                {#if attachment.uploader_name}
                  <span>آپلود شده توسط: {attachment.uploader_name}</span>
                {/if}
              </div>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2">
          <!-- Download Button -->
          <button
            onclick={() => handleDownload(attachment)}
            class="p-2 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
            title="دانلود فایل"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                    d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </button>

          <!-- Delete Button -->
          {#if canDelete}
            <button
              onclick={() => handleDelete(attachment)}
              class="p-2 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors"
              title="حذف فایل"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H8a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <!-- Summary -->
  <div class="mt-4 pt-4 border-t border-slate-200">
    <div class="flex items-center justify-between text-sm text-slate-500">
      <span>{attachments.length} فایل</span>
      <span>
        مجموع حجم: {formatFileSize(attachments.reduce((total, att) => total + att.file_size, 0))}
      </span>
    </div>
  </div>
{/if}