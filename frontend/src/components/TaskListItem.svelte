<script>
  import { formatJalaliDate } from '../lib/utils';
  
  let { task, onComplete } = $props();

  const priorityColors = {
    'Critical': 'bg-red-100 text-red-700',
    'High': 'bg-orange-100 text-orange-700',
    'Medium': 'bg-blue-100 text-blue-700',
    'Low': 'bg-gray-100 text-gray-700'
  };

  let priorityClass = $derived(priorityColors[task.priority_label] || 'bg-gray-100 text-gray-700');
  let isCompleted = $state(false);

  function handleComplete() {
    isCompleted = true;
    setTimeout(() => {
      onComplete(task.id);
    }, 1000);
  }
</script>

<div class="flex items-center p-4 hover:bg-gray-50 transition-colors {isCompleted ? 'opacity-50 grayscale' : ''}">
  <input 
    type="checkbox" 
    class="w-5 h-5 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 cursor-pointer"
    checked={isCompleted}
    onchange={handleComplete}
  />
  
  <div class="mr-4 flex-grow">
    <h5 class="text-sm font-bold text-gray-900 {isCompleted ? 'line-through' : ''}">{task.title}</h5>
    <p class="text-xs text-gray-500">{task.project_name}</p>
  </div>

  <div class="flex flex-col items-end">
    <span class="px-2 py-0.5 rounded text-[10px] font-bold mb-1 {priorityClass}">
      {task.priority_label}
    </span>
    <span class="text-[10px] text-gray-400">
      {formatJalaliDate(task.due_date, 'relative')}
    </span>
  </div>
</div>
