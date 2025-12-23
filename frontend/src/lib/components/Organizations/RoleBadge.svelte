<script lang="ts">
	import type { OrganizationMembershipRole } from '$lib/schemas/organization-membership';
	import { Badge } from '../ui/badge';

	const { role }: { role: OrganizationMembershipRole } = $props();
	const label = $derived(role.charAt(0).toUpperCase() + role.slice(1));
</script>

<Badge variant={role === 'admin' ? 'accent' : 'secondary'}>{label}</Badge>
