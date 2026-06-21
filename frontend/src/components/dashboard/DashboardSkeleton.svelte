<script>
  // Skeleton matching the bento layout
</script>

<div class="space-y-5" aria-busy="true" aria-label="در حال بارگذاری داشبورد">
  <!-- Header -->
  <div class="flex items-end justify-between mb-8">
    <div class="space-y-2">
      <div class="h-9 w-48 bg-surface-muted rounded-lg animate-pulse"></div>
      <div class="h-4 w-32 bg-surface-muted rounded animate-pulse"></div>
    </div>
    <div class="h-9 w-28 bg-surface-muted rounded-full animate-pulse"></div>
  </div>

  <!-- Stat row -->
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 md:gap-5">
    {#each Array(4) as _}
      <div class="bg-surface border border-border-subtle rounded-bento p-5 shadow-soft space-y-3">
        <div class="flex items-start justify-between">
          <div class="w-10 h-10 rounded-xl bg-surface-muted animate-pulse"></div>
          <div class="h-5 w-12 bg-surface-muted rounded-full animate-pulse"></div>
        </div>
        <div class="h-7 w-16 bg-surface-muted rounded animate-pulse"></div>
        <div class="h-3 w-24 bg-surface-muted rounded animate-pulse"></div>
      </div>
    {/each}
  </div>

  <!-- Bento row -->
  <div class="grid grid-cols-12 gap-4 md:gap-5">
    <div class="col-span-12 lg:col-span-8 bg-surface border border-border-subtle rounded-bento p-6 shadow-soft space-y-4">
      <div class="h-4 w-24 bg-surface-muted rounded animate-pulse"></div>
      <div class="h-7 w-3/4 bg-surface-muted rounded animate-pulse"></div>
      <div class="h-4 w-full bg-surface-muted rounded animate-pulse"></div>
      <div class="h-4 w-2/3 bg-surface-muted rounded animate-pulse"></div>
      <div class="flex items-center justify-between pt-4 border-t border-border-subtle">
        <div class="flex -space-x-2">
          {#each Array(3) as _}
            <div class="w-8 h-8 rounded-full bg-surface-muted animate-pulse"></div>
          {/each}
        </div>
        <div class="h-9 w-32 bg-surface-muted rounded-full animate-pulse"></div>
      </div>
    </div>
    <div class="col-span-12 lg:col-span-4 bg-surface border border-border-subtle rounded-bento p-5 shadow-soft space-y-3">
      <div class="h-4 w-28 bg-surface-muted rounded animate-pulse"></div>
      {#each Array(3) as _}
        <div class="h-14 w-full bg-surface-muted rounded-xl animate-pulse"></div>
      {/each}
    </div>
  </div>

  <!-- Bento row -->
  <div class="grid grid-cols-12 gap-4 md:gap-5">
    <div class="col-span-12 lg:col-span-8 bg-surface border border-border-subtle rounded-bento p-6 shadow-soft space-y-3">
      <div class="h-4 w-32 bg-surface-muted rounded animate-pulse"></div>
      {#each Array(4) as _}
        <div class="h-12 w-full bg-surface-muted rounded-lg animate-pulse"></div>
      {/each}
    </div>
    <div class="col-span-12 lg:col-span-4 bg-surface border border-border-subtle rounded-bento p-6 shadow-soft space-y-3">
      <div class="h-4 w-24 bg-surface-muted rounded animate-pulse"></div>
      <div class="h-20 w-full bg-surface-muted rounded-lg animate-pulse"></div>
      <div class="h-9 w-full bg-surface-muted rounded-full animate-pulse"></div>
    </div>
  </div>
</div>
