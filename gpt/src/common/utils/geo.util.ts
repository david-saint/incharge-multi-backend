export function haversineDistance(
  fromLat: number,
  fromLng: number,
  toLat: number,
  toLng: number,
  mode: 'km' | 'mi',
): number {
  const earthRadius = mode === 'mi' ? 3958.8 : 6371;
  const toRadians = (value: number) => (value * Math.PI) / 180;

  const latDistance = toRadians(toLat - fromLat);
  const lngDistance = toRadians(toLng - fromLng);

  const a =
    Math.sin(latDistance / 2) * Math.sin(latDistance / 2) +
    Math.cos(toRadians(fromLat)) *
      Math.cos(toRadians(toLat)) *
      Math.sin(lngDistance / 2) *
      Math.sin(lngDistance / 2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return earthRadius * c;
}
