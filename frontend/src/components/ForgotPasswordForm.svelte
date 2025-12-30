<script>
  import { api } from '../lib/api.js';

  // State
  let email = $state('');
  let isLoading = $state(false);
  let errorMessage = $state('');
  let successMessage = $state('');

  // Form submission
  async function handleSubmit() {
    errorMessage = '';
    successMessage = '';

    if (!email) {
      errorMessage = 'لطفاً ایمیل خود را وارد کنید';
      return;
    }

    isLoading = true;

    try {
      const response = await fetch(`${api}/auth/forgot-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        successMessage = data.data.message;
        email = ''; // Clear form
      } else {
        errorMessage = data.error?.message || 'خطا در ارسال درخواست';
      }
    } catch (error) {
      errorMessage = 'خطا در برقراری ارتباط با سرور';
      console.error('Forgot password error:', error);
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-gray-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8" dir="rtl">
  <div class="max-w-md w-full space-y-8">
    <div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        بازیابی رمز عبور
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        ایمیل خود را وارد کنید تا لینک بازیابی برای شما ارسال شود
      </p>
    </div>

    <form class="mt-8 space-y-6" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      {#if errorMessage}
        <div class="bg-red-50 border-r-4 border-red-400 p-4 rounded">
          <p class="text-sm text-red-800">{errorMessage}</p>
        </div>
      {/if}

      {#if successMessage}
        <div class="bg-green-50 border-r-4 border-green-400 p-4 rounded">
          <p class="text-sm text-green-800">{successMessage}</p>
        </div>
      {/if}

      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">
          ایمیل
        </label>
        <input
          id="email"
          type="email"
          bind:value={email}
          required
          class="mt-1 appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
          placeholder="example@email.com"
        />
      </div>

      <div>
        <button
          type="submit"
          disabled={isLoading}
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {isLoading ? 'در حال ارسال...' : 'ارسال لینک بازیابی'}
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
