<script>
  import { onMount } from 'svelte';
  import { getProjectMembers, addProjectMember, removeProjectMember, getEligibleProjectUsers, getProjectRoles } from '../lib/api/projectMembers.js';

  let { projectId, show = false } = $props();

  let members = $state([]);
  let eligibleUsers = $state([]);
  let roles = $state([]);
  let loading = $state(false);
  let error = $state(null);
  let successMessage = $state('');
  let memberToDelete = $state(null);
  let showDeleteConfirm = $state(false);

  // Form state
  let selectedUserId = $state('');
  let selectedRoleId = $state('');

  // Reset state when modal opens
  $effect(() => {
    if (show && projectId) {
      resetForm();
      loadMembers();
      loadEligibleUsers();
      loadRoles();
    }
  });

  function resetForm() {
    selectedUserId = '';
    selectedRoleId = '';
    error = null;
    successMessage = '';
  }

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
      error = 'لطفاً هم کاربر و هم نقش را انتخاب کنید';
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

  function getRoleDisplayName(roleName) {
    const roleMap = {
      'owner': 'مالک',
      'manager': 'مدیر',
      'contributor': 'مشارکت‌کننده',
      'viewer': 'بیننده'
    };
    return roleMap[roleName] || roleName;
  }

  function getRoleBadgeColor(roleName) {
    switch (roleName) {
      case 'owner': return 'bg-purple-100 text-purple-700';
      case 'manager': return 'bg-blue-100 text-blue-700';
      case 'contributor': return 'bg-green-100 text-green-700';
      case 'viewer': return 'bg-gray-100 text-gray-700';
      default: return 'bg-indigo-100 text-indigo-700';
    }
  }

  function confirmDeleteMember(member) {
    memberToDelete = member;
    showDeleteConfirm = true;
    error = null;
    successMessage = '';
  }

  function cancelDelete() {
    memberToDelete = null;
    showDeleteConfirm = false;
  }

  async function handleRemoveMember() {
    if (!memberToDelete) return;

    loading = true;
    error = null;
    successMessage = '';

    try {
      await removeProjectMember(projectId, memberToDelete.user_id);
      successMessage = 'عضو با موفقیت حذف شد';
      memberToDelete = null;
      showDeleteConfirm = false;
      
      // Reload data
      await loadMembers();
      await loadEligibleUsers();
    } catch (err) {
      error = err.message || 'خطا در حذف عضو';
    } finally {
      loading = false;
    }
  }
</script>

{#if show}
  <div class="space-y-6" dir="rtl">
    <!-- Header with member count -->
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-lg font-semibold text-gray-900">مدیریت اعضای پروژه</h3>
        <p class="text-sm text-gray-500 mt-1">{members.length} عضو در این پروژه</p>
      </div>
    </div>

    <!-- Alerts -->
    {#if error}
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl">
        {error}
      </div>
    {/if}

    {#if successMessage}
      <div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-xl">
        {successMessage}
      </div>
    {/if}

    <!-- Add Member Form -->
    <div class="bg-gray-50 rounded-xl p-4">
      <h4 class="text-sm font-medium text-gray-700 mb-4">افزودن عضو جدید</h4>
      <form onsubmit={handleAddMember} class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label for="user-select" class="block text-sm font-medium text-gray-700 mb-2">کاربر</label>
            <select
              id="user-select"
              bind:value={selectedUserId}
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white"
              disabled={loading || eligibleUsers.length === 0}
            >
              <option value="">انتخاب کاربر...</option>
              {#each eligibleUsers as user}
                <option value={user.id}>{user.username} ({user.email})</option>
              {/each}
            </select>
            {#if eligibleUsers.length === 0}
              <p class="text-xs text-gray-500 mt-1">هیچ کاربر واجد شرایطی برای افزودن وجود ندارد</p>
            {/if}
          </div>

          <div>
            <label for="role-select" class="block text-sm font-medium text-gray-700 mb-2">نقش</label>
            <select
              id="role-select"
              bind:value={selectedRoleId}
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white"
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
          class="w-full md:w-auto px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? 'در حال افزودن...' : 'افزودن عضو'}
        </button>
      </form>
    </div>

    <!-- Members List -->
    <div>
      <h4 class="text-sm font-medium text-gray-700 mb-3">اعضای فعلی</h4>
      <div class="space-y-2 max-h-80 overflow-y-auto">
        {#if loading && members.length === 0}
          <div class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          </div>
        {:else if members.length === 0}
          <p class="text-gray-500 text-center py-8 bg-gray-50 rounded-lg">هنوز عضوی به این پروژه اضافه نشده است.</p>
        {:else}
          {#each members as member}
            <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600 font-medium">
                  {member.username?.charAt(0).toUpperCase()}
                </div>
                <div>
                  <p class="font-medium text-gray-900">{member.username}</p>
                  <p class="text-sm text-gray-500">{member.email}</p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <span class="px-3 py-1 rounded-full text-sm font-medium {getRoleBadgeColor(member.role_name)}">
                  {member.role_display}
                </span>
                <button
                  onclick={() => confirmDeleteMember(member)}
                  disabled={loading}
                  class="p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-50"
                  title="حذف عضو"
                  aria-label="حذف عضو"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Delete Confirmation Dialog -->
    {#if showDeleteConfirm && memberToDelete}
      <div class="bg-red-50 border border-red-200 rounded-xl p-4">
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0 text-red-500">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <div class="flex-1">
            <h4 class="text-sm font-medium text-red-900">حذف عضو از پروژه</h4>
            <p class="text-sm text-red-700 mt-1">
              آیا مطمئن هستید که می‌خواهید <strong>{memberToDelete.username}</strong> را از پروژه حذف کنید؟
            </p>
            <div class="flex items-center gap-2 mt-3">
              <button
                onclick={handleRemoveMember}
                disabled={loading}
                class="px-4 py-2 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700 disabled:bg-red-400 disabled:cursor-not-allowed transition-colors"
              >
                {loading ? 'در حال حذف...' : 'بله، حذف شود'}
              </button>
              <button
                onclick={cancelDelete}
                disabled={loading}
                class="px-4 py-2 bg-white text-red-700 text-sm border border-red-300 rounded-lg hover:bg-red-50 disabled:opacity-50 transition-colors"
              >
                انصراف
              </button>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>
{/if}
