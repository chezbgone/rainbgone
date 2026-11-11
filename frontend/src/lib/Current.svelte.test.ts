// @vitest-environment jsdom
import { render } from '@testing-library/svelte';
import { expect, it } from 'vitest';
import Current from './Current.svelte';
import { currently, day } from '../test/fixtures';

it('updates the current icon when navigating to a different forecast', async () => {
	const time = Date.parse('2026-11-10T00:00:00Z') / 1000;
	const { container, rerender } = render(Current, {
		currently: currently(time),
		daily: { summary: 'Fixture', icon: 'clear-day', data: [day(time)] }
	});
	expect(container.querySelector('img')?.getAttribute('src')).toBe('/weather-icons/clear-day.png');
	await rerender({
		currently: currently(time, { icon: 'thunderstorm', summary: 'Thunderstorms' })
	});
	expect(container.querySelector('img')?.getAttribute('src')).toBe('/weather-icons/rain.png');
});
