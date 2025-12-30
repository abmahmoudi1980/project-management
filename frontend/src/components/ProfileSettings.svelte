<script>
  import { authStore } from '../stores/authStore.js';

  let activeTab = $state('profile');
  let profileForm = $state({
    username: $authStore.user?.username || '',
    email: $authStore.user?.email || ''
  });
  let passwordForm = $state({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  let isLoading = $state(false);
  let error = $state('');
  let success = $state('');

  async function updateProfile() {
    if (!profileForm.username || !profileForm.email) {
      error = 'نام کاربری و ایمیل نمی‌توانند خالی باشند';
      return;
    }

    isLoading = true;
    error = '';
    success = '';

    try {
      const response = await fetch('/api/auth/me', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(profileForm)
      });

      const data = await response.json();

      if (response.ok) {
        success = 'پروفایل با موفقیت بروزرسانی شد';
        // Update the store
        authStore.update(state => ({ ...state, user: { ...state.user, ...profileForm } }));
      } else {
        error = data.error?.message || 'خطا در بروزرسانی پروفایل';
      }
    } catch (err) {
      error = 'خطا در ارتباط با سرور';
    } finally {
      isLoading = false;
    }
  }

  async function changePassword() {
    if (!passwordForm.currentPassword || !passwordForm.newPassword) {
      error = 'تمام فیلدها را پر کنید';
      return;
    }

    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      error = 'رمز عبور جدید و تکرار آن مطابقت ندارند';
      return;
    }

    if (passwordForm.newPassword.length < 8) {
      error = 'رمز عبور جدید باید حداقل 8 کاراکتر باشد';
      return;
    }

    isLoading = true;
    error = '';
    success = '';

    try {
      const response = await fetch('/api/auth/me/password', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          current_password: passwordForm.currentPassword,
          new_password: passwordForm.newPassword
        })
      });

      const data = await response.json();

      if (response.ok) {
        success = 'رمز عبور با موفقیت تغییر یافت';
        passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };
      } else {
        error = data.error?.message || 'خطا در تغییر رمز عبور';
      }
    } catch (err) {
      error = 'خطا در ارتباط با سرور';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-md">
  <h2 class="text-2xl font-bold mb-6 text-center">تنظیمات پروفایل</h2>

  <!-- Tab Navigation -->
  <div class="flex mb-6 bg-gray-100 rounded-lg p-1">
    <button
      class="flex-1 py-2 px-4 rounded-md text-center transition-colors {activeTab === 'profile' ? 'bg-white shadow-sm' : 'text-gray-600'}"
      onclick={() => { activeTab = 'profile'; error = ''; success = ''; }}
    >
      ویرایش پروفایل
    </button>
    <button
      class="flex-1 py-2 px-4 rounded-md text-center transition-colors {activeTab === 'password' ? 'bg-white shadow-sm' : 'text-gray-600'}"
      onclick={() => { activeTab = 'password'; error = ''; success = ''; }}
    >
      تغییر رمز عبور
    </button>
  </div>

  <!-- Error/Success Messages -->
  {#if error}
    <div class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
      {error}
    </div>
  {/if}

  {#if success}
    <div class="mb-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded">
      {success}
    </div>
  {/if}

  <!-- Profile Form -->
  {#if activeTab === 'profile'}
    <form on:submit|preventDefault={updateProfile} class="space-y-4">
      <div>
        <label for="username" class="block text-sm font-medium text-gray-700 mb-1">
          نام کاربری
        </label>
        <input
          id="username"
          type="text"
          bind:value={profileForm.username}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="نام کاربری خود را وارد کنید"
          required
        />
      </div>

      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
          ایمیل
        </label>
        <input
          id="email"
          type="email"
          bind:value={profileForm.email}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="ایمیل خود را وارد کنید"
          required
        />
      </div>

      <button
        type="submit"
        disabled={isLoading}
        class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
      >
        {#if isLoading}
          در حال بروزرسانی...
        {:else}
          بروزرسانی پروفایل
        {/if}
      </button>
    </form>
  {/if}

  <!-- Password Form -->
  {#if activeTab === 'password'}
    <form on:submit|preventDefault={changePassword} class="space-y-4">
      <div>
        <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">
          رمز عبور فعلی
        </label>
        <input
          id="currentPassword"
          type="password"
          bind:value={passwordForm.currentPassword}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="رمز عبور فعلی خود را وارد کنید"
          required
        />
      </div>

      <div>
        <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">
          رمز عبور جدید
        </label>
        <input
          id="newPassword"
          type="password"
          bind:value={passwordForm.newPassword}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="رمز عبور جدید را وارد کنید"
          required
        />
      </div>

      <div>
        <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">
          تکرار رمز عبور جدید
        </label>
        <input
          id="confirmPassword"
          type="password"
          bind:value={passwordForm.confirmPassword}
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="رمز عبور جدید را تکرار کنید"
          required
        />
      </div>

      <button
        type="submit"
        disabled={isLoading}
        class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
      >
        {#if isLoading}
          در حال تغییر رمز عبور...
        {:else}
          تغییر رمز عبور
        {/if}
      </button>
    </form>
  {/if}
</div>