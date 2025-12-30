<script>
  import { authStore } from '../stores/authStore.js';
  
  // State using Svelte 5 runes
  let email = $state('');
  let password = $state('');
  let error = $state('');
  let isLoading = $state(false);
  
  // Derived validation
  let isValid = $derived(
    email.trim().includes('@') &&
    password.length >= 1
  );
  
  // Form submission
  async function handleSubmit() {
    error = '';
    isLoading = true;
    
    const result = await authStore.login(
      email.trim().toLowerCase(),
      password
    );
    
    isLoading = false;
    
    if (!result.success) {
      error = result.error;
    }
  }
</script>

<div class="max-w-md mx-auto mt-8 p-6 bg-white rounded-lg shadow-md" dir="rtl">
  <h2 class="text-2xl font-bold mb-6 text-center text-gray-800">ورود</h2>
  
  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
    <!-- Email -->
    <div class="mb-4">
      <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
        ایمیل
      </label>
      <input
        type="email"
        id="email"
        bind:value={email}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        required
      />
    </div>
    
    <!-- Password -->
    <div class="mb-6">
      <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
        رمز عبور
      </label>
      <input
        type="password"
        id="password"
        bind:value={password}
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
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
      class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
    >
      {isLoading ? 'در حال ورود...' : 'ورود'}
    </button>
  </form>
</div>
