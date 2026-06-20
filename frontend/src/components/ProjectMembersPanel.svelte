<script>
  import { onMount } from 'svelte';
  import { getProjectMembers, addProjectMember, getEligibleProjectUsers, getProjectRoles } from '../lib/api/projectMembers.js';

  let { projectId } = $props();

  let members = $state([]);
  let eligibleUsers = $state([]);
  let roles = $state([]);
  let loading = $state(false);
  let error = $state(null);
  let successMessage = $state('');

  // Form state
  let selectedUserId = $state('');
  let selectedRoleId = $state('');

  async function loadMembers() {
    if (!projectId) return;
    
    loading = true;
    error = null;
    try {
      const result = await getProjectMembers(projectId);
      members = result || [];
    } catch (err) {
      error = err.message;
      members = [];
    } finally {
      loading = false;
    }
  }

  async function loadEligibleUsers() {
    if (!projectId) return;
    
    try {
      const result = await getEligibleProjectUsers(projectId);
      eligibleUsers = result || [];
    } catch (err) {
      console.error('Failed to load eligible users:', err);
      eligibleUsers = [];
    }
  }

  async function loadRoles() {
    try {
      const result = await getProjectRoles();
      roles = result || [];
    } catch (err) {
      console.error('Failed to load roles:', err);
      roles = [];
    }
  }

  async function handleAddMember(event) {
    event.preventDefault();
    
    if (!selectedUserId || !selectedRoleId) {
      error = 'Please select both a user and a role';
      return;
    }

    loading = true;
    error = null;
    successMessage = '';

    try {
      await addProjectMember(projectId, {
        user_id: selectedUserId,
        project_role_id: selectedRoleId
      });
      
      successMessage = 'عضو با موفقیت اضافه شد';
      selectedUserId = '';
      selectedRoleId = '';
      
      // Reload data
      await loadMembers();
      await loadEligibleUsers();
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    loadMembers();
    loadEligibleUsers();
    loadRoles();
  });

  $effect(() => {
    if (projectId) {
      loadMembers();
      loadEligibleUsers();
    }
  });
</script>

<div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6" dir="rtl">
  <h2 class="text-xl font-bold text-slate-900 mb-4">اعضای پروژه</h2>

  {#if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl mb-4">
      {error}
    </div>
  {/if}

  {#if successMessage}
    <div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-xl mb-4">
      {successMessage}
    </div>
  {/if}

  <!-- Add Member Form -->
  <form onsubmit={handleAddMember} class="mb-6 space-y-4">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">کاربر</label>
        <select
          bind:value={selectedUserId}
          class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
          disabled={loading || eligibleUsers.length === 0}
        >
          <option value="">انتخاب کاربر...</option>
          {#each eligibleUsers as user}
            <option value={user.id}>{user.username} ({user.email})</option>
          {/each}
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">نقش</label>
        <select
          bind:value={selectedRoleId}
          class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-brand-500"
          disabled={loading || roles.length === 0}
        >
          <option value="">انتخاب نقش...</option>
          {#each roles as role}
            <option value={role.id}>{role.display_name}</option>
          {/each}
        </select>
      </div>
    </div>

    <button
      type="submit"
      disabled={loading || !selectedUserId || !selectedRoleId}
      class="w-full md:w-auto px-6 py-2 bg-brand-600 text-white rounded-lg hover:bg-brand-700 disabled:bg-slate-400 disabled:cursor-not-allowed transition-colors"
    >
      {loading ? 'در حال افزودن...' : 'افزودن عضو'}
    </button>
  </form>

  <!-- Members List -->
  <div class="space-y-2">
    {#if loading && members.length === 0}
      <div class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-600"></div>
      </div>
    {:else if members.length === 0}
      <p class="text-slate-500 text-center py-8">هنوز عضوی به این پروژه اضافه نشده است.</p>
    {:else}
      {#each members as member}
        <div class="flex items-center justify-between p-4 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-brand-100 flex items-center justify-center text-brand-600 font-medium">
              {member.username?.charAt(0).toUpperCase()}
            </div>
            <div>
              <p class="font-medium text-slate-900">{member.username}</p>
              <p class="text-sm text-slate-500">{member.email}</p>
            </div>
          </div>
          <span class="px-3 py-1 bg-brand-100 text-brand-700 rounded-full text-sm">
            {member.role_display}
          </span>
        </div>
      {/each}
    {/if}
  </div>
</div>
