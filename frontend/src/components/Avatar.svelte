<script>
	let { user, size = 'md' } = $props();

	const sizeClasses = {
		sm: 'w-7 h-7 text-xs',
		md: 'w-8 h-8 text-sm',
		lg: 'w-10 h-10 text-base'
	};

	const colors = [
		'bg-indigo-500',
		'bg-purple-500',
		'bg-blue-500',
		'bg-pink-500',
		'bg-orange-500',
		'bg-green-500',
		'bg-red-500',
		'bg-teal-500'
	];

	function getInitials(fullName) {
		const parts = fullName.trim().split(' ');
		if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
		return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
	}

	function getColorClass(userId) {
		const hash = userId.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
		return colors[hash % colors.length];
	}

	const initials = $derived(getInitials(user.full_name));
	const colorClass = $derived(getColorClass(user.id));
</script>

<div
	class="rounded-full {sizeClasses[size]} {colorClass} flex items-center justify-center text-white font-bold border-2 border-white"
	title={user.full_name}
>
	{initials}
</div>
