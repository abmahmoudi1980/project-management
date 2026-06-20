<script>
  import { authStore } from '../stores/authStore.js';
  import Button from './ui/Button.svelte';
  import Input from './ui/Input.svelte';

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let isLoading = $state(false);

  let isValid = $derived(
    email.trim().includes('@') &&
    password.length >= 1
  );

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

<div class="max-w-md mx-auto mt-4 sm:mt-8 bg-white rounded-lg shadow-md" dir="rtl">
  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="p-4 sm:p-6 space-y-4">
    <h2 class="text-xl sm:text-2xl font-bold text-center text-slate-800">ورود</h2>

    {#if error}
      <div class="p-3 bg-danger-50 border border-danger-200 text-danger-700 rounded-md text-sm" role="alert">
        {error}
      </div>
    {/if}

    <Input
      type="email"
      label="ایمیل"
      bind:value={email}
      autocomplete="email"
      required
    />

    <Input
      type="password"
      label="رمز عبور"
      bind:value={password}
      autocomplete="current-password"
      required
    />

    <Button type="submit" disabled={!isValid} loading={isLoading} fullWidth>
      {isLoading ? 'در حال ورود...' : 'ورود'}
    </Button>
  </form>
</div>
