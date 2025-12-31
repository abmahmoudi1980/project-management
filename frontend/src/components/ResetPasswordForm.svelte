<script>
  import { api } from '../lib/api.js';

  // Props - get token from URL
  let { token = '' } = $props();

  // State
  let newPassword = $state('');
  let confirmPassword = $state('');
  let isLoading = $state(false);
  let errorMessage = $state('');
  let successMessage = $state('');

  // Derived - client-side validation
  let passwordsMatch = $derived(newPassword === confirmPassword);
  let isPasswordStrong = $derived(
    newPassword.length >= 8 &&
    /[A-Z]/.test(newPassword) &&
    /[a-z]/.test(newPassword) &&
    /[0-9]/.test(newPassword)
  );

  // Form submission
  async function handleSubmit() {
    errorMessage = '';
    successMessage = '';

    // Validation
    if (!newPassword || !confirmPassword) {
      errorMessage = 'لطفاً تمام فیلدها را پر کنید';
      return;
    }

    if (!passwordsMatch) {
      errorMessage = 'رمزهای عبور مطابقت ندارند';
      return;
    }

    if (!isPasswordStrong) {
      errorMessage = 'رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد';
      return;
    }

    isLoading = true;

    try {
      const response = await fetch(`${api}/auth/reset-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ 
          token: token,
          new_password: newPassword 
        }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        successMessage = data.data.message;
        newPassword = '';
        confirmPassword = '';
        
        // Redirect to login after 2 seconds
        setTimeout(() => {
          window.location.hash = '#/login';
        }, 2000);
      } else {
        errorMessage = data.error?.message || 'خطا در تغییر رمز عبور';
      }
    } catch (error) {
      errorMessage = 'خطا در برقراری ارتباط با سرور';
      console.error('Reset password error:', error);
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-gray-100 flex items-center justify-center py-8 sm:py-12 px-4 sm:px-6 lg:px-8" dir="rtl">
  <div class="max-w-md w-full space-y-6 sm:space-y-8">
    <div>
      <h2 class="mt-4 sm:mt-6 text-center text-2xl sm:text-3xl font-extrabold text-gray-900">
        تنظیم رمز عبور جدید
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        رمز عبور جدید خود را وارد کنید
      </p>
    </div>

    <form class="mt-6 sm:mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      {#if errorMessage}
        <div class="bg-red-50 border-r-4 border-red-400 p-4 rounded">
          <p class="text-sm text-red-800">{errorMessage}</p>
        </div>
      {/if}

      {#if successMessage}
        <div class="bg-green-50 border-r-4 border-green-400 p-4 rounded">
          <p class="text-sm text-green-800">{successMessage}</p>
          <p class="text-xs text-green-600 mt-2">در حال انتقال به صفحه ورود...</p>
        </div>
      {/if}

      <div>
        <label for="newPassword" class="block text-sm font-medium text-gray-700">
          رمز عبور جدید
        </label>
        <input
          id="newPassword"
          type="password"
          bind:value={newPassword}
          required
          class="mt-1 appearance-none rounded-md relative block w-full px-3 py-3 min-h-[44px] border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          placeholder="حداقل 8 کاراکتر"
        />
        {#if newPassword && !isPasswordStrong}
          <p class="mt-1 text-xs text-red-600">
            رمز عبور باید حداقل 8 کاراکتر و شامل حروف بزرگ، کوچک و اعداد باشد
          </p>
        {/if}
      </div>

      <div>
        <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
          تکرار رمز عبور
        </label>
        <input
          id="confirmPassword"
          type="password"
          bind:value={confirmPassword}
          required
          class="mt-1 appearance-none rounded-md relative block w-full px-3 py-3 min-h-[44px] border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          placeholder="تکرار رمز عبور"
        />
        {#if confirmPassword && !passwordsMatch}
          <p class="mt-1 text-xs text-red-600">
            رمزهای عبور مطابقت ندارند
          </p>
        {/if}
      </div>

      <div>
        <button
          type="submit"
          disabled={isLoading || !passwordsMatch || !isPasswordStrong}
          class="group relative w-full flex justify-center py-2.5 px-4 min-h-[44px] border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {isLoading ? 'در حال تغییر رمز عبور...' : 'تغییر رمز عبور'}
        </button>
      </div>

      <div class="text-center">
        <a href="#/login" class="text-sm text-blue-600 hover:text-blue-500">
          بازگشت به صفحه ورود
        </a>
      </div>
    </form>
  </div>
</div>

<style>
  input {
    text-align: right;
    direction: rtl;
  }
</style>
