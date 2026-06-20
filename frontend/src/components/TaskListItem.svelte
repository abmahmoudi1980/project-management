<script>
  import { formatJalaliDate } from '../lib/utils';
  
  let { task, onComplete } = $props();

  const priorityColors = {
    'Critical': 'bg-red-100 text-red-700',
    'High': 'bg-orange-100 text-orange-700',
    'Medium': 'bg-brand-100 text-brand-700',
    'Low': 'bg-slate-100 text-slate-700'
  };

  let priorityClass = $derived(priorityColors[task.priority_label] || 'bg-slate-100 text-slate-700');
  let isCompleted = $state(false);

  function handleComplete() {
    isCompleted = true;
    setTimeout(() => {
      onComplete(task.id);
    }, 1000);
  }
</script>

<div class="flex items-center p-4 hover:bg-slate-50 transition-colors {isCompleted ? 'opacity-50 grayscale' : ''}">
  <input
    type="checkbox"
    class="w-5 h-5 rounded border-slate-300 text-brand-600 focus:ring-brand-500 cursor-pointer"
    checked={isCompleted}
    onchange={handleComplete}
    aria-label={isCompleted ? `لغو تکمیل ${task.title}` : `تکمیل ${task.title}`}
  />
  
  <div class="mr-4 flex-grow">
    <h5 class="text-sm font-bold text-slate-900 {isCompleted ? 'line-through' : ''}">{task.title}</h5>
    <p class="text-xs text-slate-500">{task.project_name}</p>
  </div>

  <div class="flex flex-col items-end">
    <span class="px-2 py-0.5 rounded text-[10px] font-bold mb-1 {priorityClass}">
      {task.priority_label}
    </span>
    <span class="text-[10px] text-slate-400">
      {formatJalaliDate(task.due_date, 'relative')}
    </span>
  </div>
</div>
