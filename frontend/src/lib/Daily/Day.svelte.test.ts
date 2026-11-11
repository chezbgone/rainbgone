// @vitest-environment jsdom
import { render, within } from '@testing-library/svelte';
import { expect, it } from 'vitest';
import Day from './Day.svelte';
import { day, hour } from '../../test/fixtures';

it('updates icons and temperature-bar geometry when an existing day receives new props', async () => {
	const time = Date.parse('2026-11-10T00:00:00Z') / 1000;
	const { container, rerender } = render(Day, {
		daily: day(time, { temperatureMin: 10, temperatureMax: 30 }),
		hourly: [hour(time), hour(time + 3600)],
		minTemp: 0,
		maxTemp: 100,
		today: false,
		lat: 0,
		lng: 0,
		timezone: 'UTC'
	});
	const summary = container.querySelector('summary')!;
	expect(summary.querySelector('img')?.getAttribute('src')).toBe('/weather-icons/clear-day.png');
	expect((within(summary).getByText('10°') as HTMLElement).style.left).toBe('10%');

	await rerender({
		daily: day(time, { icon: 'thunderstorm', temperatureMin: 40, temperatureMax: 80 }),
		minTemp: 20,
		maxTemp: 120
	});
	expect(summary.querySelector('img')?.getAttribute('src')).toBe('/weather-icons/rain.png');
	expect((within(summary).getByText('40°') as HTMLElement).style.left).toBe('20%');
	expect((within(summary).getByText('80°') as HTMLElement).style.left).toBe('60%');
	const bar = summary.querySelector<HTMLElement>('.rounded-full')!;
	expect(bar.style.marginLeft).toBe('20%');
	expect(bar.style.width).toBe('40%');
});
