<script>
  /**
   * Quick search box for real-time text search
   * Searches task title and description as user types
   */
  let { searchText = $bindable('') } = $props();
</script>

<div class="search-box">
  <input
    type="text"
    bind:value={searchText}
    placeholder="جستجو در عنوان و توضیح..."
    class="w-full px-4 py-3 border border-slate-300 rounded-lg text-sm
      placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500
      focus:border-transparent transition-colors"
  />
</div>

<style>
  .search-box {
    display: flex;
    width: 100%;
  }
</style>
