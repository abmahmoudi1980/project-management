<script>
  import { authStore } from '../stores/authStore.js';
  
  // State using Svelte 5 runes
  let username = $state('');
  let email = $state('');
  let password = $state('');
  let passwordConfirmation = $state('');
  let error = $state('');
  let isLoading = $state(false);
  
  // Derived validation
  let isValid = $derived(
    username.trim().length >= 3 &&
    email.trim().includes('@') &&
    password.length >= 8 &&
    password === passwordConfirmation
  );
  
  // Form submission
  async function handleSubmit() {
    error = '';
    
    // Client-side validation
    if (password !== passwordConfirmation) {
      error = 'رمز عبور و تکرار آن مطابقت ندارند';
      return;
    }
    
    isLoading = true;
    
    const result = await authStore.register(
      username.trim(),
      email.trim().toLowerCase(),
      password,
      passwordConfirmation
    );
    
    isLoading = false;
    
    if (!result.success) {
      error = result.error;
    }
  }
</script>

<div class="max-w-md mx-auto mt-4 sm:mt-8 p-4 sm:p-6 bg-white rounded-lg shadow-md" dir="rtl">
  <h2 class="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-center text-gray-800">ثبت‌نام</h2>

  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
    <!-- Username -->
    <div class="mb-4">
      <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
        نام کاربری
      </label>
      <input
        type="text"
        id="username"
        bind:value={username}
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        required
        minlength="3"
        maxlength="50"
      />
    </div>

    <!-- Email -->
    <div class="mb-4">
      <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
        ایمیل
      </label>
      <input
        type="email"
        id="email"
        bind:value={email}
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        required
      />
    </div>

    <!-- Password -->
    <div class="mb-4">
      <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
        رمز عبور
      </label>
      <input
        type="password"
        id="password"
        bind:value={password}
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        required
        minlength="8"
      />
      <p class="text-xs text-gray-500 mt-1">
        حداقل 8 کاراکتر، شامل حروف بزرگ، کوچک و اعداد
      </p>
    </div>

    <!-- Password Confirmation -->
    <div class="mb-6">
      <label for="passwordConfirmation" class="block text-sm font-medium text-gray-700 mb-2">
        تکرار رمز عبور
      </label>
      <input
        type="password"
        id="passwordConfirmation"
        bind:value={passwordConfirmation}
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        required
      />
    </div>

    <!-- Error Message -->
    {#if error}
      <div class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-md text-sm">
        {error}
      </div>
    {/if}

    <!-- Submit Button -->
    <button
      type="submit"
      disabled={!isValid || isLoading}
      class="w-full min-h-[44px] bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium"
    >
      {isLoading ? 'در حال ثبت‌نام...' : 'ثبت‌نام'}
    </button>
  </form>
</div>
