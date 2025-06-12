// src/app/utils/chart.utils.ts

/**
 * Resamples time series data at fixed intervals using linear interpolation
 * @param data Original data points: {x: number (timestamp), y: number}[]
 * @param interval Sampling interval in milliseconds
 * @returns Resampled data with interpolated points
 */
export function resampleData(
  data: { x: number; y: number }[],
  interval: number
): { x: number; y: number }[] {
  if (!data.length) return [];

  // Sort by timestamp ascending
  const sorted = [...data].sort((a, b) => a.x - b.x);
  const first = sorted[0].x;
  const last = sorted[sorted.length - 1].x;
  const pointsByTime = new Map<number, number>(
    sorted.map(d => [d.x, d.y])
  );

  const resampled: { x: number; y: number }[] = [];
  for (let t = first; t <= last; t += interval) {
    resampled.push({
      x: t,
      y: pointsByTime.get(t) ?? 0
    });
  }

  return resampled;
}