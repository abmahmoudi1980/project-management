<script>
  let { user, size = 'md' } = $props();

  const sizes = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-10 h-10 text-sm',
    lg: 'w-12 h-12 text-base'
  };

  const colors = [
    'bg-indigo-100 text-indigo-700',
    'bg-purple-100 text-purple-700',
    'bg-blue-100 text-blue-700',
    'bg-pink-100 text-pink-700',
    'bg-orange-100 text-orange-700',
    'bg-green-100 text-green-700',
    'bg-red-100 text-red-700',
    'bg-teal-100 text-teal-700'
  ];

  function getInitials(name) {
    if (!name) return '?';
    const parts = name.trim().split(' ');
    if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }

  function getColor(id) {
    if (!id) return colors[0];
    const hash = id.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return colors[hash % colors.length];
  }

  let initials = $derived(getInitials(user?.full_name || user?.username || user?.email));
  let colorClass = $derived(getColor(user?.id));
</script>

<div class="{sizes[size]} {colorClass} rounded-full flex items-center justify-center font-medium overflow-hidden flex-shrink-0" title={user?.full_name || user?.username}>
  {initials}
</div>
