<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';

  // State
  let users = $state([]);
  let isLoading = $state(true);
  let errorMessage = $state('');
  let successMessage = $state('');
  let pagination = $state({ page: 1, limit: 20, total: 0, total_pages: 0 });
  
  // Confirmation dialog state
  let showConfirmDialog = $state(false);
  let confirmAction = $state(null);
  let confirmMessage = $state('');

  onMount(async () => {
    await loadUsers();
  });

  async function loadUsers() {
    isLoading = true;
    errorMessage = '';

    try {
      const data = await api.users.getAll(pagination.page, pagination.limit);
      users = data.data.users || [];
      pagination = data.data.pagination;
    } catch (error) {
      errorMessage = 'خطا در برقراری ارتباط با سرور';
      console.error('Load users error:', error);
    } finally {
      isLoading = false;
    }
  }

  async function changeUserRole(userId, newRole) {
    try {
      await api.users.updateRole(userId, newRole);
      successMessage = 'نقش کاربر با موفقیت تغییر یافت';
      setTimeout(() => (successMessage = ''), 3000);
      await loadUsers();
    } catch (error) {
      errorMessage = 'خطا در برقراری ارتباط با سرور';
      console.error('Change role error:', error);
    }
  }

  async function toggleUserActivation(userId, isActive) {
    try {
      await api.users.updateActivation(userId, !isActive);
      successMessage = isActive
        ? 'کاربر با موفقیت غیرفعال شد'
        : 'کاربر با موفقیت فعال شد';
      setTimeout(() => (successMessage = ''), 3000);
      await loadUsers();
    } catch (error) {
      errorMessage = 'خطا در برقراری ارتباط با سرور';
      console.error('Toggle activation error:', error);
    }
  }

  function confirmRoleChange(userId, currentRole) {
    const newRole = currentRole === 'admin' ? 'user' : 'admin';
    const roleLabel = newRole === 'admin' ? 'ادمین' : 'کاربر عادی';
    
    confirmMessage = `آیا از تغییر نقش این کاربر به "${roleLabel}" اطمینان دارید؟`;
    confirmAction = () => changeUserRole(userId, newRole);
    showConfirmDialog = true;
  }

  function confirmActivationToggle(userId, username, isActive) {
    const action = isActive ? 'غیرفعال' : 'فعال';
    confirmMessage = `آیا از ${action} کردن کاربر "${username}" اطمینان دارید؟`;
    confirmAction = () => toggleUserActivation(userId, isActive);
    showConfirmDialog = true;
  }

  function executeConfirmAction() {
    if (confirmAction) {
      confirmAction();
    }
    closeConfirmDialog();
  }

  function closeConfirmDialog() {
    showConfirmDialog = false;
    confirmAction = null;
    confirmMessage = '';
  }

  function nextPage() {
    if (pagination.page < pagination.total_pages) {
      pagination.page++;
      loadUsers();
    }
  }

  function previousPage() {
    if (pagination.page > 1) {
      pagination.page--;
      loadUsers();
    }
  }

  function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('fa-IR');
  }
</script>

<div class="p-6" dir="rtl">
  <div class="mb-6">
    <h2 class="text-2xl font-bold text-gray-900">مدیریت کاربران</h2>
    <p class="text-sm text-gray-600 mt-1">مشاهده و مدیریت کاربران سیستم</p>
  </div>

  {#if errorMessage}
    <div class="bg-red-50 border-r-4 border-red-400 p-4 rounded mb-4">
      <p class="text-sm text-red-800">{errorMessage}</p>
      <button
        onclick={() => (errorMessage = '')}
        class="text-red-600 hover:text-red-800 text-xs mt-2"
      >
        بستن
      </button>
    </div>
  {/if}

  {#if successMessage}
    <div class="bg-green-50 border-r-4 border-green-400 p-4 rounded mb-4">
      <p class="text-sm text-green-800">{successMessage}</p>
    </div>
  {/if}

  {#if isLoading}
    <div class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="mr-3 text-gray-600">در حال بارگذاری...</span>
    </div>
  {:else}
    <!-- Desktop Table -->
    <div class="hidden md:block bg-white shadow-md rounded-lg overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              نام کاربری
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              ایمیل
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              نقش
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              وضعیت
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              تاریخ عضویت
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">
              عملیات
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#each users as user (user.id)}
            <tr>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {user.username}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                {user.email}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  class="px-2 py-1 text-xs font-medium rounded {user.role === 'admin'
                    ? 'bg-purple-100 text-purple-800'
                    : 'bg-gray-100 text-gray-800'}"
                >
                  {user.role === 'admin' ? 'ادمین' : 'کاربر عادی'}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  class="px-2 py-1 text-xs font-medium rounded {user.is_active
                    ? 'bg-green-100 text-green-800'
                    : 'bg-red-100 text-red-800'}"
                >
                  {user.is_active ? 'فعال' : 'غیرفعال'}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                {formatDate(user.created_at)}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm space-x-2 space-x-reverse">
                <button
                  onclick={() => confirmRoleChange(user.id, user.role)}
                  class="text-blue-600 hover:text-blue-800 font-medium"
                >
                  تغییر نقش
                </button>
                <span class="text-gray-300">|</span>
                <button
                  onclick={() =>
                    confirmActivationToggle(user.id, user.username, user.is_active)}
                  class="{user.is_active
                    ? 'text-red-600 hover:text-red-800'
                    : 'text-green-600 hover:text-green-800'} font-medium"
                >
                  {user.is_active ? 'غیرفعال کردن' : 'فعال کردن'}
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Mobile Card Layout -->
    <div class="md:hidden space-y-3">
      {#each users as user (user.id)}
        <div class="bg-white rounded-lg shadow-md p-4">
          <div class="flex items-start justify-between mb-3">
            <div>
              <h3 class="text-base font-semibold text-gray-900">
                {user.username}
              </h3>
              <p class="text-sm text-gray-600 mt-0.5">{user.email}</p>
            </div>
            <div class="flex flex-col gap-2">
              <span
                class="px-2 py-1 text-xs font-medium rounded {user.role === 'admin'
                  ? 'bg-purple-100 text-purple-800'
                  : 'bg-gray-100 text-gray-800'}"
              >
                {user.role === 'admin' ? 'ادمین' : 'کاربر عادی'}
              </span>
              <span
                class="px-2 py-1 text-xs font-medium rounded {user.is_active
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800'}"
              >
                {user.is_active ? 'فعال' : 'غیرفعال'}
              </span>
            </div>
          </div>
          <div class="text-xs text-gray-500 mb-3">
            تاریخ عضویت: {formatDate(user.created_at)}
          </div>
          <div class="flex flex-col gap-2">
            <button
              onclick={() => confirmRoleChange(user.id, user.role)}
              class="w-full px-4 py-2.5 min-h-[44px] bg-blue-50 text-blue-700 font-medium rounded-lg hover:bg-blue-100 transition-colors text-sm"
            >
              تغییر نقش
            </button>
            <button
              onclick={() => confirmActivationToggle(user.id, user.username, user.is_active)}
              class="w-full px-4 py-2.5 min-h-[44px] {user.is_active
                ? 'bg-red-50 text-red-700 hover:bg-red-100'
                : 'bg-green-50 text-green-700 hover:bg-green-100'} font-medium rounded-lg transition-colors text-sm"
            >
              {user.is_active ? 'غیرفعال کردن' : 'فعال کردن'}
            </button>
          </div>
        </div>
      {/each}
    </div>

    {#if users.length === 0}
      <div class="text-center py-12 text-gray-500 px-4">
        کاربری یافت نشد
      </div>
    {/if}

    {#if pagination.total_pages > 1}
      <div class="mt-4 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
        <div class="text-sm text-gray-600 text-center md:text-right">
          صفحه {pagination.page} از {pagination.total_pages} (کل: {pagination.total} کاربر)
        </div>
        <div class="flex space-x-2 space-x-reverse justify-center md:justify-start">
          <button
            onclick={previousPage}
            disabled={pagination.page === 1}
            class="px-4 py-2.5 min-h-[44px] min-w-[44px] text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            قبلی
          </button>
          <button
            onclick={nextPage}
            disabled={pagination.page === pagination.total_pages}
            class="px-4 py-2.5 min-h-[44px] min-w-[44px] text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            بعدی
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>

<!-- Confirmation Dialog -->
{#if showConfirmDialog}
  <div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center z-50 p-4" dir="rtl">
    <div class="bg-white rounded-lg shadow-xl p-4 sm:p-6 max-w-md w-full">
      <h3 class="text-lg font-medium text-gray-900 mb-4">تأیید عملیات</h3>
      <p class="text-gray-600 mb-6">{confirmMessage}</p>
      <div class="flex flex-col sm:flex-row justify-end gap-3">
        <button
          onclick={closeConfirmDialog}
          class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
        >
          انصراف
        </button>
        <button
          onclick={executeConfirmAction}
          class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
        >
          تأیید
        </button>
      </div>
    </div>
  </div>
{/if}
