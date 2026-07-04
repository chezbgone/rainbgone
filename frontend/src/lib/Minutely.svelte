<script lang="ts">
	import type { Forecast } from '$lib/Forecast/types';

	interface Props {
		minutely: Forecast['minutely'];
	}

	interface Point {
		intensity: number;
		error: number;
	}

	let { minutely }: Props = $props();
	let canvas = $state<HTMLCanvasElement>();
	let frameId = 0;
	let wobbleStart = 0;

	const width = 350;
	const height = 120;

	const hasPrecip = $derived(minutely.data.some((minute) => minute.precipIntensity > 0));
	const minutes = $derived(minutely.data.slice(0, 61));
	const clamp = (value: number, min = 0, max = 1) => Math.min(Math.max(value, min), max);

	const transformIntensity = (rawIntensity: number) => {
		const mediumIntensity = 1 / 3;
		const transformed = 1 - Math.exp(-2.209389806 * Math.sqrt(rawIntensity));

		if (transformed <= 0.5) {
			return transformed * (mediumIntensity / 0.5);
		}
		return mediumIntensity + (transformed - 0.5) * ((1 - mediumIntensity) / 0.5);
	};

	const points: Point[] = $derived(
		minutes.map((minute) => {
			const intensity = transformIntensity(minute.precipIntensity);

			return {
				intensity,
				error:
					minute.precipIntensity > 0
						? 0.0025 * Math.min(10, minute.precipIntensityError / minute.precipIntensity)
						: 0
			};
		})
	);

	const wobbleTable = [
		26, 25, 25, 24, 22, 20, 18, 17, 16, 16, 17, 20, 23, 26, 29, 31, 33, 33, 32, 30, 27, 24, 21, 19,
		18, 17, 17, 19, 21, 25, 28, 30, 32, 33, 33, 32, 30, 28, 27, 25, 24, 24, 25, 26, 27, 29, 29, 30,
		29, 28, 27, 26, 25, 24, 23, 21, 20, 18, 17, 17, 17, 16, 16, 16, 15, 16, 17, 19, 22, 25, 28, 30,
		32, 34, 36, 37, 39, 40, 40, 39, 37, 35, 32, 30, 29, 27, 26, 25, 23, 22, 21, 22, 24, 28, 32, 35,
		38, 39, 38, 37, 36, 35, 35, 35, 35, 35, 35, 34, 33, 32, 32, 33, 34, 37, 39, 40, 41, 40, 38, 35,
		32, 29, 27, 26, 26, 26, 27, 28, 28, 28, 27, 25, 23, 21, 19, 17, 16, 15, 14, 14, 13, 13, 13, 12,
		12, 11, 10, 10, 10, 11, 13, 15, 17, 20, 22, 25, 28, 32, 34, 35, 35, 33, 29, 24, 18, 13, 9, 6, 4,
		2, 0, 0, 1, 3, 6, 8, 9, 10, 9, 8, 7, 7, 8, 9, 10, 11, 11, 10, 7, 4, 2, 0, 0, 1, 4, 7, 12, 16,
		19, 21, 22, 22, 21, 21, 23, 25, 28, 32, 35, 37, 38, 39, 40, 40, 39, 39, 37, 35, 33, 30, 26, 23,
		19, 15, 12, 10, 8, 7, 8, 9, 11, 15, 18, 22, 26, 30, 33, 35, 36, 36, 36, 36, 36, 38, 39, 40, 41,
		41, 39, 37, 35, 32, 30, 28, 27, 26
	];

	const wobbleOffset = (x: number, time: number) => {
		const negativeIndex = ((x - time) * wobbleTable.length) | 0;
		const positiveIndex = ((x + time) * wobbleTable.length) | 0;

		return (
			8 *
			(wobbleTable[(positiveIndex + wobbleTable.length) % wobbleTable.length] -
				wobbleTable[(negativeIndex + wobbleTable.length) % wobbleTable.length])
		);
	};

	const drawFrame = () => {
		if (!canvas) return;

		const context = canvas.getContext('2d');
		if (!context) return;

		const devicePixelRatio = window.devicePixelRatio || 1;
		const time = (0.0001 * (Date.now() - wobbleStart)) % 1;

		if (canvas.width !== width * devicePixelRatio || canvas.height !== height * devicePixelRatio) {
			canvas.width = width * devicePixelRatio;
			canvas.height = height * devicePixelRatio;
		}

		context.setTransform(devicePixelRatio, 0, 0, devicePixelRatio, 0, 0);
		context.clearRect(0, 0, width, height);

		context.fillStyle = '#80a5d6';
		context.beginPath();

		for (let index = points.length; index--; ) {
			const point = points[index];
			const normalizedX = points.length > 0 ? index / points.length : 0;
			const x = points.length > 1 ? (index * width) / (points.length - 1) : 0;
			const wobbleIntensity = clamp(point.error, 0, 0.00015);
			const intensity = clamp(point.intensity + wobbleIntensity * wobbleOffset(normalizedX, time));

			context.lineTo(x, height * (1 - intensity));
		}

		context.lineTo(-10, height);
		context.lineTo(width, height);
		context.fill();

		frameId = requestAnimationFrame(drawFrame);
	};

	$effect(() => {
		if (!hasPrecip || !canvas) {
			cancelAnimationFrame(frameId);
			return;
		}

		wobbleStart = Date.now();
		frameId = requestAnimationFrame(drawFrame);

		return () => {
			cancelAnimationFrame(frameId);
		};
	});
</script>

{#if hasPrecip}
	<div class="mx-auto pt-2 text-center">
		<div class="text-lg font-light">Next Hour: {minutely.summary}</div>

		<div class="inline-flex items-start pt-4">
			<div class="w-[350px]">
				<div class="relative h-[120px] w-[350px]">
					<canvas bind:this={canvas} {width} {height} class="block h-[120px] w-[350px]"></canvas>
					<div class="pointer-events-none absolute inset-x-0 top-0 h-full">
						{#each Array(3) as _}
							<div
								class="h-[33.333333%] border-r border-dashed border-neutral-400 not-last:border-b"
							></div>
						{/each}
					</div>
				</div>

				<!-- ticks -->
				<div class="flex justify-between">
					{#each Array(7) as _}
						<span class="h-2 border-l-2 border-neutral-400"></span>
					{/each}
				</div>

				<!-- labels -->
				<div class="flex h-4 justify-between pb-8">
					{#each Array(7) as _, i}
						<div class="relative">
							{#if i % 2 === 1}
								<div class="absolute -translate-x-[calc(50%-2px)] text-sm whitespace-nowrap">
									{10 * i} min
								</div>
							{/if}
						</div>
					{/each}
				</div>
			</div>

			<div class="ml-2 flex h-[120px] flex-col justify-around text-left text-xs font-light">
				<div>Heavy</div>
				<div>Medium</div>
				<div>Light</div>
			</div>
		</div>
	</div>
{/if}
