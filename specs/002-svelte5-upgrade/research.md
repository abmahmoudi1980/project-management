# Research: Svelte 5 Migration

**Feature**: Upgrade Frontend to Svelte 5  
**Branch**: `002-svelte5-upgrade`  
**Date**: 2025-12-30

## Overview

Svelte 5 introduces a major paradigm shift with the **runes system**, replacing the previous reactivity model while maintaining backward compatibility where possible. This document consolidates research findings to guide the migration from Svelte 4.2.0 to Svelte 5.

## Key Changes in Svelte 5

### 1. Runes System (Core Reactivity)

Svelte 5 replaces the compile-time magic with explicit runes for managing reactive state.

#### `$state` - Reactive State

**Purpose**: Declare reactive state variables.

**Svelte 4 Pattern**:
```svelte
<script>
  let count = 0;
</script>
```

**Svelte 5 Pattern**:
```svelte
<script>
  let count = $state(0);
</script>
```

**Key Differences**:
- Explicit declaration of reactive variables
- Works inside and outside components
- Can be used in classes and regular JavaScript
- Mutating nested properties is reactive by default

**Deep Reactivity Example**:
```svelte
<script>
  let todos = $state([
    { id: 1, text: 'Buy milk', done: false }
  ]);
  
  // This is reactive - no need for assignment
  function toggle(id) {
    const todo = todos.find(t => t.id === id);
    todo.done = !todo.done;
  }
</script>
```

#### `$derived` - Computed Values

**Purpose**: Derive values from reactive state (replaces `$:` reactive declarations).

**Svelte 4 Pattern**:
```svelte
<script>
  let count = 0;
  $: doubled = count * 2;
  $: quadrupled = doubled * 2;
</script>
```

**Svelte 5 Pattern**:
```svelte
<script>
  let count = $state(0);
  let doubled = $derived(count * 2);
  let quadrupled = $derived(doubled * 2);
</script>
```

**Key Benefits**:
- Explicit dependency tracking
- No ordering issues (unlike `$:` statements)
- Can be used anywhere, not just in components
- Better TypeScript support

#### `$effect` - Side Effects

**Purpose**: Run side effects in response to state changes (replaces `$:` blocks).

**Svelte 4 Pattern**:
```svelte
<script>
  let count = 0;
  
  $: {
    console.log('count changed:', count);
    document.title = `Count: ${count}`;
  }
</script>
```

**Svelte 5 Pattern**:
```svelte
<script>
  let count = $state(0);
  
  $effect(() => {
    console.log('count changed:', count);
    document.title = `Count: ${count}`;
  });
</script>
```

**Cleanup Pattern**:
```svelte
<script>
  $effect(() => {
    const interval = setInterval(() => {
      count++;
    }, 1000);
    
    return () => clearInterval(interval);
  });
</script>
```

**Key Rules**:
- Automatic dependency tracking (uses what you access)
- Return function for cleanup
- Don't use for derived state (use `$derived` instead)
- Runs after component mounts and when dependencies change

#### `$props` - Component Props

**Purpose**: Declare component properties (replaces `export let`).

**Svelte 4 Pattern**:
```svelte
<script>
  export let title = 'Default Title';
  export let count = 0;
  export let onClick = () => {};
</script>
```

**Svelte 5 Pattern**:
```svelte
<script>
  let { title = 'Default Title', count = 0, onClick = () => {} } = $props();
</script>
```

**Destructuring with Rest**:
```svelte
<script>
  let { title, ...restProps } = $props();
</script>

<div {...restProps}>
  <h1>{title}</h1>
</div>
```

**Key Benefits**:
- Standard JavaScript destructuring syntax
- Better TypeScript inference
- Easier to pass through props with spread
- More explicit and readable

### 2. Stores Compatibility

**Good News**: Svelte stores (writable, readable, derived) remain fully compatible in Svelte 5.

**No Changes Needed For**:
- Store creation: `writable()`, `readable()`, `derived()`
- Store subscriptions: `$storeName` syntax still works
- Store methods: `set()`, `update()`, `subscribe()`

**Example - Stores Work As-Is**:
```javascript
// projectStore.js - NO CHANGES NEEDED
import { writable } from 'svelte/store';

function createProjectStore() {
  const { subscribe, set, update } = writable([]);
  
  return {
    subscribe,
    load: async () => { /* ... */ },
    create: async (project) => { /* ... */ }
  };
}

export const projects = createProjectStore();
```

```svelte
<!-- Component using store - NO CHANGES NEEDED -->
<script>
  import { projects } from './stores/projectStore';
</script>

{#each $projects as project}
  <div>{project.title}</div>
{/each}
```

**When to Consider Migration**: If you want to use runes for better performance or to access state outside components, you can migrate stores to runes-based state, but it's **not required**.

### 3. Event Handling Changes

**Component Events**: Event dispatching remains the same in Svelte 5.

**Svelte 4/5 - No Change**:
```svelte
<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  
  function handleClick() {
    dispatch('select', { id: 123 });
  }
</script>

<button on:click={handleClick}>Select</button>
```

**Usage in Parent - No Change**:
```svelte
<MyComponent on:select={handleSelect} />
```

### 4. Lifecycle Hooks

**onMount, onDestroy, beforeUpdate, afterUpdate** - All remain available and work the same way.

**No Changes Needed**:
```svelte
<script>
  import { onMount } from 'svelte';
  
  onMount(async () => {
    // Initialization code
    return () => {
      // Cleanup code
    };
  });
</script>
```

**Note**: For new code, prefer `$effect` for side effects, but existing `onMount` usage doesn't need to change.

### 5. Two-Way Binding

**bind:** directive works the same way with `$state`.

**Svelte 4**:
```svelte
<script>
  let value = '';
</script>
<input bind:value />
```

**Svelte 5**:
```svelte
<script>
  let value = $state('');
</script>
<input bind:value />
```

**No special changes needed** - bindings work identically.

### 6. Conditional Rendering & Loops

**{#if}, {#each}, {#await}** - All work exactly the same way.

**No Changes Needed**:
```svelte
{#if condition}
  <p>True</p>
{:else}
  <p>False</p>
{/if}

{#each items as item}
  <div>{item.name}</div>
{/each}

{#await promise}
  <p>Loading...</p>
{:then value}
  <p>Loaded: {value}</p>
{:catch error}
  <p>Error: {error}</p>
{/await}
```

## Migration Strategy for This Project

### Components to Migrate

1. **App.svelte** (Priority: High)
   - Migrate: `selectedProject` state → `$state`
   - Migrate: `onMount` (keep as-is or convert to `$effect`)
   - Verify: Store subscriptions still work

2. **ProjectForm.svelte, TaskForm.svelte, TimeLogForm.svelte** (Priority: High)
   - Migrate: Props from `export let` → `$props()`
   - Migrate: Form state variables → `$state`
   - Verify: Two-way bindings work
   - Verify: Form submission logic

3. **ProjectList.svelte, TaskList.svelte** (Priority: Medium)
   - Migrate: Props → `$props()`
   - Migrate: Local state → `$state`
   - Verify: Event dispatching works
   - Verify: Store subscriptions work

4. **Modal.svelte** (Priority: Medium)
   - Migrate: Props (show, onClose) → `$props()`
   - Verify: Modal visibility toggling

5. **JalaliDatePicker.svelte** (Priority: High - Third-party dependency)
   - Migrate: Props → `$props()`
   - Migrate: Calendar state → `$state`
   - Test thoroughly: Date selection, Jalali conversion

### Stores (No Migration Needed)

The existing store implementations are compatible:
- `projectStore.js` - Keep as-is
- `taskStore.js` - Keep as-is
- `timeLogStore.js` - Keep as-is

### Build Configuration

**package.json updates**:
```json
{
  "devDependencies": {
    "svelte": "^5.0.0",
    "@sveltejs/vite-plugin-svelte": "^5.0.0",
    "vite": "^5.0.0"
  }
}
```

**vite.config.js** - Likely no changes needed, but verify plugin options.

## Common Patterns Reference

### Pattern 1: Simple Component State

**Before**:
```svelte
<script>
  let isOpen = false;
  
  function toggle() {
    isOpen = !isOpen;
  }
</script>
```

**After**:
```svelte
<script>
  let isOpen = $state(false);
  
  function toggle() {
    isOpen = !isOpen;
  }
</script>
```

### Pattern 2: Props with Defaults

**Before**:
```svelte
<script>
  export let title = 'Untitled';
  export let description = '';
  export let status = 'active';
</script>
```

**After**:
```svelte
<script>
  let { 
    title = 'Untitled', 
    description = '', 
    status = 'active' 
  } = $props();
</script>
```

### Pattern 3: Derived Values

**Before**:
```svelte
<script>
  let firstName = '';
  let lastName = '';
  
  $: fullName = `${firstName} ${lastName}`.trim();
  $: initials = `${firstName[0] || ''}${lastName[0] || ''}`.toUpperCase();
</script>
```

**After**:
```svelte
<script>
  let firstName = $state('');
  let lastName = $state('');
  
  let fullName = $derived(`${firstName} ${lastName}`.trim());
  let initials = $derived(`${firstName[0] || ''}${lastName[0] || ''}`.toUpperCase());
</script>
```

### Pattern 4: Side Effects

**Before**:
```svelte
<script>
  let count = 0;
  
  $: {
    if (count > 10) {
      console.warn('Count exceeded 10');
    }
  }
</script>
```

**After**:
```svelte
<script>
  let count = $state(0);
  
  $effect(() => {
    if (count > 10) {
      console.warn('Count exceeded 10');
    }
  });
</script>
```

### Pattern 5: Component with Events

**Before** (works in both):
```svelte
<script>
  import { createEventDispatcher } from 'svelte';
  export let item;
  
  const dispatch = createEventDispatcher();
  
  function handleClick() {
    dispatch('select', item);
  }
</script>

<button on:click={handleClick}>Select</button>
```

**After** (same code works):
```svelte
<script>
  import { createEventDispatcher } from 'svelte';
  let { item } = $props();
  
  const dispatch = createEventDispatcher();
  
  function handleClick() {
    dispatch('select', item);
  }
</script>

<button on:click={handleClick}>Select</button>
```

## Testing Checklist

After migrating each component:

- [ ] Component renders without errors
- [ ] Props are received correctly
- [ ] State updates trigger re-renders
- [ ] Derived values compute correctly
- [ ] Effects run at appropriate times
- [ ] Events dispatch correctly
- [ ] Two-way bindings work
- [ ] Store subscriptions work
- [ ] No console errors or warnings

## Resources

- **Official Svelte 5 Docs**: https://svelte.dev/docs/svelte/overview
- **Migration Guide**: https://svelte.dev/docs/svelte/v5-migration-guide
- **Runes Tutorial**: https://learn.svelte.dev/tutorial/runes
- **What's New in Svelte 5**: https://svelte.dev/blog/svelte-5-is-alive

## Decision Log

### Decision 1: Keep Stores As-Is

**Rationale**: Svelte stores are fully compatible with Svelte 5. Migrating them to runes would be unnecessary work with minimal benefit for this project's scale.

**Alternatives Considered**: 
- Migrate stores to runes-based state management
- Use global `$state` variables

**Rejected Because**: Current store pattern works well, is well-tested, and provides good encapsulation of logic.

### Decision 2: Migrate All Components to Runes

**Rationale**: While Svelte 5 has backward compatibility mode, using runes provides better developer experience, clearer code, and positions the codebase for future Svelte development.

**Alternatives Considered**:
- Keep components in compatibility mode
- Mix old and new syntax

**Rejected Because**: Mixed syntax would be confusing and harder to maintain. Full migration ensures consistency.

### Decision 3: Preserve Lifecycle Hooks

**Rationale**: `onMount` and other lifecycle hooks still work and are more explicit about lifecycle timing than generic `$effect`.

**When to Use What**:
- `onMount` - For initialization that should happen once
- `$effect` - For reactive side effects based on state changes
- Both are valid in Svelte 5

## Known Issues & Workarounds

### Issue 1: Vite Plugin Version

**Problem**: `@sveltejs/vite-plugin-svelte` version 3.x is for Svelte 4, version 5.x is for Svelte 5.

**Solution**: Must upgrade plugin to version 5.x alongside Svelte 5.

### Issue 2: Deep Reactivity Differences

**Problem**: Svelte 5's deep reactivity might behave slightly differently than Svelte 4.

**Solution**: Test object/array mutations carefully. If issues arise, use explicit `$state` cloning patterns.

### Issue 3: TypeScript Support

**Problem**: This project doesn't use TypeScript, but Svelte 5's runes are designed with TypeScript in mind.

**Impact**: None - runes work perfectly well in plain JavaScript.

**Future Consideration**: Svelte 5 would make TypeScript migration easier if desired later.

## Performance Expectations

Based on Svelte 5 announcements and benchmarks:

- **Bundle Size**: 30-50% smaller than Svelte 4 (due to improved compiler)
- **Runtime Performance**: Faster reactivity updates (runes are more efficient)
- **Build Time**: Similar or slightly faster
- **HMR**: Same or better hot module replacement speed

**For This Project**: Expect bundle size to decrease from ~95KB to ~60-80KB (estimated).
